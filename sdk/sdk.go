package sdk

import (
	"bytes"
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/opensourceways/community-robot-lib/utils"

	"github.com/opensourceways/repo-file-cache/models"
)

func NewSDK(endpoint string, maxRetries int) *SDK {
	slash := "/"
	if !strings.HasSuffix(endpoint, slash) {
		endpoint += slash
	}

	return &SDK{
		hc:              &utils.HttpClient{MaxRetries: maxRetries},
		endpoint:        endpoint,
		summaryEndpoint: endpoint + "%s?summary=true",
	}
}

type SDK struct {
	hc              *utils.HttpClient
	endpoint        string
	summaryEndpoint string
}

func (cli *SDK) SaveFiles(opt *models.FileUpdateOption) error {
	payload, err := utils.JsonMarshal(opt)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, cli.endpoint, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	return cli.ForwardTo(req, nil)
}

func (cli *SDK) GetFiles(b models.Branch, fileName string, summary bool) (models.FilesInfo, error) {
	endpoint := ""
	if summary {
		endpoint = fmt.Sprintf(cli.summaryEndpoint, genFilePath(b, fileName))
	} else {
		endpoint = cli.endpoint + genFilePath(b, fileName)
	}

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return models.FilesInfo{}, err
	}

	var v struct {
		Data models.FilesInfo `json:"data"`
	}

	if err = cli.ForwardTo(req, &v); err != nil {
		return models.FilesInfo{}, err
	}

	return v.Data, nil
}

func (cli *SDK) DeleteFiles(b models.Branch, fileName string, opt *models.FileDeleteOption) error {
	payload, err := utils.JsonMarshal(opt)
	if err != nil {
		return err
	}

	endpoint := cli.endpoint + genFilePath(b, fileName)

	req, err := http.NewRequest(http.MethodDelete, endpoint, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	return cli.ForwardTo(req, nil)
}

func (cli *SDK) ForwardTo(req *http.Request, jsonResp interface{}) error {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "repo-file-cache-sdk")

	return cli.hc.ForwardTo(req, jsonResp)
}

func genFilePath(b models.Branch, fileName string) string {
	return path.Join(b.Platform, b.Org, b.Repo, b.Branch, fileName)
}
