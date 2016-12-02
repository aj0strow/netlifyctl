package porcelain

import (
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"

	"encoding/hex"

	"io/ioutil"

	"bytes"

	"github.com/netlify/open-api/go/models"
	"github.com/netlify/open-api/go/plumbing/operations"
	"github.com/netlify/open-api/go/porcelain/context"
)

const (
	preProcessingTimeout = time.Minute * 5
)

type uploadError struct {
	err   error
	mutex *sync.Mutex
}

type file struct {
	Name   string
	SHA1   string
	Buffer []byte
}

type deployFiles struct {
	Files  map[string]*file
	Sums   map[string]string
	Hashed map[string]*file
}

func newDeployFiles() *deployFiles {
	return &deployFiles{
		Files:  make(map[string]*file),
		Sums:   make(map[string]string),
		Hashed: make(map[string]*file),
	}
}

func (d *deployFiles) Add(p string, f *file) {
	d.Files[p] = f
	d.Sums[p] = f.SHA1
	d.Hashed[f.SHA1] = f
}

func (n *Netlify) overCommitted(d *deployFiles) bool {
	return len(d.Files) > n.syncFileLimit
}

// DeploySite creates a new deploy for a site given a directory in the filesystem.
// It uploads the necessary files that changed between deploys.
func (n *Netlify) DeploySite(ctx context.Context, siteID, dir string) (*models.Deploy, error) {
	f, err := os.Stat(dir)
	if err != nil {
		return nil, err
	}
	if !f.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", dir)
	}

	files, err := walk(dir)
	if err != nil {
		return nil, err
	}

	return n.createDeploy(ctx, siteID, files)
}

func (n *Netlify) createDeploy(ctx context.Context, siteID string, files *deployFiles) (*models.Deploy, error) {
	deployFiles := &models.DeployFiles{
		Files: files.Sums,
		Async: n.overCommitted(files),
	}
	l := context.GetLogger(ctx)
	l.WithFields(logrus.Fields{
		"site_id":      siteID,
		"deploy_files": len(files.Sums),
	}).Debug("Starting to deploy files")
	authInfo := context.GetAuthInfo(ctx)

	params := operations.NewCreateSiteDeployParams().WithSiteID(siteID).WithDeploy(deployFiles)
	resp, err := n.Operations.CreateSiteDeploy(params, authInfo)
	if err != nil {
		return nil, err
	}

	deploy := resp.Payload
	if n.overCommitted(files) {
		var err error
		deploy, err = n.WaitUntilDeployReady(ctx, deploy)
		if err != nil {
			return nil, err
		}
	}

	l.Debugf("Site and deploy created, uploading %d files necessary", len(deploy.Required))
	if err := n.uploadFiles(ctx, deploy, files); err != nil {
		return nil, err
	}

	return deploy, nil
}

func (n *Netlify) WaitUntilDeployReady(ctx context.Context, d *models.Deploy) (*models.Deploy, error) {
	authInfo := context.GetAuthInfo(ctx)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	params := operations.NewGetSiteDeployParams().WithSiteID(d.SiteID).WithDeployID(d.ID)
	start := time.Now()
	for t := range ticker.C {
		resp, err := n.Operations.GetSiteDeploy(params, authInfo)
		if err != nil {
			time.Sleep(3 * time.Second)
			continue
		}
		context.GetLogger(ctx).WithFields(logrus.Fields{
			"deploy_id": d.ID,
			"state":     resp.Payload.State,
		}).Debug("Waiting until deploy ready")

		if resp.Payload.State == "prepared" || resp.Payload.State == "ready" {
			return resp.Payload, nil
		}

		if resp.Payload.State == "error" {
			return nil, fmt.Errorf("Error: preprocessing deploy failed")
		}

		if t.Sub(start) > preProcessingTimeout {
			return nil, fmt.Errorf("Error: preprocessing deploy timed out")
		}
	}

	return d, nil
}

func (n *Netlify) uploadFiles(ctx context.Context, d *models.Deploy, files *deployFiles) error {
	sharedErr := &uploadError{err: nil, mutex: &sync.Mutex{}}
	sem := make(chan int, n.uploadLimit)
	wg := &sync.WaitGroup{}

	for _, sha := range d.Required {
		if file, exist := files.Hashed[sha]; exist {
			sem <- 1
			wg.Add(1)

			go n.uploadFile(ctx, d, file, wg, sem, sharedErr)
		}
	}

	wg.Wait()

	return sharedErr.err
}

func (n *Netlify) uploadFile(ctx context.Context, d *models.Deploy, f *file, wg *sync.WaitGroup, sem chan int, sharedErr *uploadError) {
	defer func() {
		wg.Done()
		<-sem
	}()

	l := context.GetLogger(ctx).WithFields(logrus.Fields{
		"deploy_id": d.ID,
		"file_path": f.Name,
	})

	sharedErr.mutex.Lock()
	if sharedErr.err != nil {
		sharedErr.mutex.Unlock()
		l.Debug("Skipping because of previous error")
		return
	}
	sharedErr.mutex.Unlock()

	authInfo := context.GetAuthInfo(ctx)

	l.WithField("file_sum", f.SHA1).Debug("Uploading file")

	//
	//b := backoff.NewExponentialBackOff()
	//b.MaxElapsedTime = 2 * time.Minute
	//
	//err := backoff.Retry(func() error {
	//	sharedErr.mutex.Lock()
	//
	//	if sharedErr.err != nil {
	//		sharedErr.mutex.Unlock()
	//		return fmt.Errorf("Upload cancelled: %s", f.Name)
	//	}
	//	sharedErr.mutex.Unlock()
	//
	fmt.Println("about to set: " + string(f.Buffer))

	params := operations.NewUploadDeployFileParams().WithDeployID(d.ID).WithPath(f.Name)
	params = params.WithFileBody(ioutil.NopCloser(bytes.NewBuffer(f.Buffer)))
	_, err := n.Operations.UploadDeployFile(params, authInfo)

	if err != nil {
		l.WithError(err).Errorf("Failed to upload %s", f.Name)
	}
	//
	//	return err
	//}, b)

	if err != nil {
		sharedErr.mutex.Lock()
		sharedErr.err = err
		sharedErr.mutex.Unlock()
	}
}

func walk(dir string) (*deployFiles, error) {
	files := newDeployFiles()

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && info.Mode().IsRegular() {
			rel, err := filepath.Rel(dir, path)
			if err != nil {
				return err
			}

			if ignoreFile(rel) {
				return nil
			}

			//o, err := os.Open(path)
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			//hash :=
			//buf := new(bytes.Buffer)
			//
			//m := io.MultiWriter(hash, buf)
			//
			//if _, err := io.Copy(m, o); err != nil {
			//	return err
			//}

			fmt.Println("contents are '" + string(data) + "'")
			file := &file{
				Name:   rel,
				Buffer: data,
				SHA1:   hex.EncodeToString(sha1.New().Sum(data)),
			}

			files.Add(rel, file)
		}

		return nil
	})

	return files, err
}

func ignoreFile(rel string) bool {
	if strings.HasPrefix(rel, ".") || strings.Contains(rel, "/.") || strings.HasPrefix(rel, "__MACOS") {
		return !strings.HasPrefix(rel, ".well-known/")
	}
	return false
}
