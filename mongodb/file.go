package mongodb

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/opensourceways/repo-file-cache/dbmodels"
)

func dirToKey(s string) (string, dbmodels.IDBError) {
	if s == dbmodels.FileCurrentDir {
		return fieldFilesItem + dbmodels.FilePathSeparator, nil
	}

	if strings.Contains(s, dbmodels.FileCurrentDir) {
		err := dbmodels.NewDBError(
			dbmodels.ErrInvalidFilePath,
			fmt.Errorf("file path contains '%s'", dbmodels.FileCurrentDir),
		)
		return "", err
	}

	return fieldFilesItem + s, nil
}

func keyToDir(s string) string {
	if s == dbmodels.FilePathSeparator {
		return ""
	}
	return s
}

func docFilterOfFiles(b *dbmodels.Branch, fileName string) bson.M {
	branch := fmt.Sprintf("%s:%s:%s:%s", b.Platform, b.Org, b.Repo, b.Branch)
	return bson.M{
		fieldBranch: branch,
		fieldName:   fileName,
	}
}

func toDocOfFiles(branchSHA string, files []dbmodels.File) (bson.M, dbmodels.IDBError) {
	m := make(bson.M, len(files))
	for i := range files {
		item := &files[i]

		key, err := dirToKey(item.Dir())
		if err != nil {
			return nil, err
		}

		m[key] = bson.M{
			fieldKey:   item.SHA,
			fieldValue: item.Content,
		}
	}

	m[fieldSHA] = branchSHA
	return m, nil
}

func (cl *client) UpdateFiles(branch *dbmodels.Branch, branchSHA string, files []dbmodels.File) dbmodels.IDBError {
	if len(files) == 0 {
		return nil
	}

	var err dbmodels.IDBError

	withContext(func(ctx context.Context) {
		docFilter := docFilterOfFiles(branch, files[0].Name())

		var doc bson.M
		doc, err = toDocOfFiles(branchSHA, files)
		if err == nil {
			err = cl.insertDocIfNotExist(ctx, cl.filesCollection, docFilter, doc)
		}
	})

	return err
}

func (cl *client) GetFiles(branch *dbmodels.Branch, fileName string) (string, []dbmodels.File, dbmodels.IDBError) {
	var v dFiles
	var err dbmodels.IDBError

	withContext(func(ctx context.Context) {
		err = cl.getDoc(
			ctx, cl.filesCollection,
			docFilterOfFiles(branch, fileName),
			bson.M{
				fieldSHA:   1,
				fieldFiles: 1,
			},
			&v,
		)
	})

	if err != nil {
		return "", nil, err
	}

	r := make([]dbmodels.File, 0, len(v.Files))
	for key, info := range v.Files {
		r = append(r, dbmodels.File{
			Path:    filepath.Join(keyToDir(key), fileName),
			SHA:     info.SHA,
			Content: info.Content,
		})
	}
	return v.BranchSHA, r, nil
}
