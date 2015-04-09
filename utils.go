package main

import "github.com/logicminds/seed"
import "os"
import "io"
import "path/filepath"
import "path"
import "strings"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// syncs all the repos defined in the shmock file and then recursively copies
// all the shmock template files to the specified shmocks directory
func sync_shmock_file(shmocks_dir string, shmock_file string) {
	repo_dir := os.TempDir()
	// the seed library will sync all the repos to a temporary directory
	repo_dirs := seed.Sync_Seed_File(repo_dir, shmock_file)
	for _, repo_dir := range repo_dirs {
		// then we need to copy over all the files to the shmocks directory
		_ = filepath.Walk(repo_dir, func(path string, info os.FileInfo, err error) error {
			if strings.HasSuffix(path, ".tmpl") {
				new_path := strings.Replace(path, repo_dir, shmocks_dir, 1)
				err = CopyFile(path, new_path)
				check(err)
			}
			return nil
		})
	}

}

// Copy the file from source to destination, while ensuring the base directory exists first
func CopyFile(source string, dest string) (err error) {
	filepath := path.Dir(dest)
	// make the parent directory if it doesn not exist
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		os.MkdirAll(filepath, 0744)
	}
	sf, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sf.Close()
	df, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer df.Close()
	_, err = io.Copy(df, sf)
	if err == nil {
		si, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, si.Mode())
		}
	}
	return
}
