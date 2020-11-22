package gitlab

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/httpclient"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {

	// The actual provider
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GITLAB_TOKEN", nil),
				Description: descriptions["token"],
			},
			"base_url": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("GITLAB_BASE_URL", ""),
				Description:  descriptions["base_url"],
				ValidateFunc: validateApiURLVersion,
			},
			"cacert_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: descriptions["cacert_file"],
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: descriptions["insecure"],
			},
			"client_cert": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: descriptions["client_cert"],
			},
			"client_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: descriptions["client_key"],
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"apgitlab_group":            dataSourceGitlabGroup(),
			"apgitlab_group_membership": dataSourceGitlabGroupMembership(),
			"apgitlab_project":          dataSourceGitlabProject(),
			"apgitlab_projects":         dataSourceGitlabProjects(),
			"apgitlab_user":             dataSourceGitlabUser(),
			"apgitlab_users":            dataSourceGitlabUsers(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"apgitlab_branch_protection":          resourceGitlabBranchProtection(),
			"apgitlab_tag_protection":             resourceGitlabTagProtection(),
			"apgitlab_group":                      resourceGitlabGroup(),
			"apgitlab_project":                    resourceGitlabProject(),
			"apgitlab_label":                      resourceGitlabLabel(),
			"apgitlab_group_label":                resourceGitlabGroupLabel(),
			"apgitlab_pipeline_schedule":          resourceGitlabPipelineSchedule(),
			"apgitlab_pipeline_schedule_variable": resourceGitlabPipelineScheduleVariable(),
			"apgitlab_pipeline_trigger":           resourceGitlabPipelineTrigger(),
			"apgitlab_project_hook":               resourceGitlabProjectHook(),
			"apgitlab_deploy_key":                 resourceGitlabDeployKey(),
			"apgitlab_deploy_key_enable":          resourceGitlabDeployEnableKey(),
			"apgitlab_deploy_token":               resourceGitlabDeployToken(),
			"apgitlab_user":                       resourceGitlabUser(),
			"apgitlab_project_membership":         resourceGitlabProjectMembership(),
			"apgitlab_group_membership":           resourceGitlabGroupMembership(),
			"apgitlab_project_variable":           resourceGitlabProjectVariable(),
			"apgitlab_group_variable":             resourceGitlabGroupVariable(),
			"apgitlab_project_cluster":            resourceGitlabProjectCluster(),
			"apgitlab_service_slack":              resourceGitlabServiceSlack(),
			"apgitlab_service_jira":               resourceGitlabServiceJira(),
			"apgitlab_service_github":             resourceGitlabServiceGithub(),
			"apgitlab_service_pipelines_email":    resourceGitlabServicePipelinesEmail(),
			"apgitlab_project_share_group":        resourceGitlabProjectShareGroup(),
			"apgitlab_group_cluster":              resourceGitlabGroupCluster(),
			"apgitlab_group_ldap_link":            resourceGitlabGroupLdapLink(),
			"apgitlab_instance_cluster":           resourceGitlabInstanceCluster(),
			"apgitlab_project_mirror":             resourceGitlabProjectMirror(),
			"apgitlab_project_level_mr_approvals": resourceGitlabProjectLevelMRApprovals(),
			"apgitlab_project_approval_rule":      resourceGitlabProjectApprovalRule(),
		},
	}

	provider.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		return providerConfigure(provider, d)
	}

	return provider
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"token": "The OAuth token used to connect to GitLab.",

		"base_url": "The GitLab Base API URL",

		"cacert_file": "A file containing the ca certificate to use in case ssl certificate is not from a standard chain",

		"insecure": "Disable SSL verification of API calls",

		"client_cert": "File path to client certificate when GitLab instance is behind company proxy. File  must contain PEM encoded data.",

		"client_key": "File path to client key when GitLab instance is behind company proxy. File must contain PEM encoded data.",
	}
}

func providerConfigure(p *schema.Provider, d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Token:      d.Get("token").(string),
		BaseURL:    d.Get("base_url").(string),
		CACertFile: d.Get("cacert_file").(string),
		Insecure:   d.Get("insecure").(bool),
		ClientCert: d.Get("client_cert").(string),
		ClientKey:  d.Get("client_key").(string),
	}

	client, err := config.Client()
	if err != nil {
		return nil, err
	}

	// NOTE: httpclient.TerraformUserAgent is deprecated and removed in Terraform SDK v2
	// After upgrading the SDK to v2 replace with p.UserAgent("terraform-provider-gitlab")
	client.UserAgent = httpclient.TerraformUserAgent(p.TerraformVersion) + " terraform-provider-gitlab"

	return client, err
}

func validateApiURLVersion(value interface{}, key string) (ws []string, es []error) {
	v := value.(string)
	if strings.HasSuffix(v, "/api/v3") || strings.HasSuffix(v, "/api/v3/") {
		es = append(es, fmt.Errorf("terraform-provider-gitlab does not support v3 api; please upgrade to /api/v4 in %s", v))
	}
	return
}
