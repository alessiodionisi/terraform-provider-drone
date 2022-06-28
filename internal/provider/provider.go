package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/drone/drone-go/drone"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/oauth2"
)

func init() {
	schema.DescriptionKind = schema.StringMarkdown

	schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
		desc := s.Description
		if s.Default != nil {
			desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
		}
		if s.Deprecated != "" {
			desc += " " + s.Deprecated
		}
		return strings.TrimSpace(desc)
	}
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"server": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The Drone servers url, It must be provided, but can also be sourced from the `DRONE_SERVER` environment variable.",
					DefaultFunc: schema.EnvDefaultFunc("DRONE_SERVER", nil),
				},
				"token": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The Drone servers api token, It must be provided, but can also be sourced from the `DRONE_TOKEN` environment variable.",
					DefaultFunc: schema.EnvDefaultFunc("DRONE_TOKEN", nil),
				},
			},
			ResourcesMap: map[string]*schema.Resource{
				"drone_repo":      resourceRepo(),
				"drone_secret":    resourceSecret(),
				"drone_orgsecret": resourceOrgSecret(),
				"drone_user":      resourceUser(),
				"drone_cron":      resourceCron(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {
		server := data.Get("server").(string)
		token := data.Get("token").(string)

		oauth2Conf := &oauth2.Config{}
		httpClient := oauth2Conf.Client(ctx, &oauth2.Token{
			AccessToken: token,
		})

		droneClient := drone.NewClient(server, httpClient)

		if _, err := droneClient.Self(); err != nil {
			return nil, diag.Errorf("drone client failed: %s", err)
		}

		return droneClient, nil
	}
}

// func configure(ctx context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {
// 	config := new(oauth2.Config)

// 	tlsConfig := &tls.Config{
// 		InsecureSkipVerify: false,
// 	}

// 	auther := config.Client(
// 		oauth2.NoContext,
// 		&oauth2.Token{
// 			AccessToken: data.Get("token").(string),
// 		},
// 	)

// 	trans, _ := auther.Transport.(*oauth2.Transport)
// 	trans.Base = &http.Transport{
// 		TLSClientConfig: tlsConfig,
// 		Proxy:           http.ProxyFromEnvironment,
// 	}

// 	client := drone.NewClient(data.Get("server").(string), auther)

// 	if _, err := client.Self(); err != nil {
// 		return nil, diag.Errorf("drone client failed: %s", err)
// 	}

// 	return client, diag.Diagnostics{}
// }
