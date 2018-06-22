package azurerm

import (
	"encoding/base64"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2016-09-01/web"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceArmAppServiceCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppServiceCertificateCreateUpdate,
		Read:   resourceAppServiceCertificateRead,
		Update: resourceAppServiceCertificateCreateUpdate,
		Delete: resourceAppServiceCertificateDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_group_name": resourceGroupNameSchema(),
			"location":            locationSchema(),
			"base_64_encoded_pfx_file": {
				Type:     schema.TypeString,
				Required: true,
			},
			"pfx_password": {
				Type:     schema.TypeString,
				Required: true,
			},
			"thumbprint": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAppServiceCertificateCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).appServicesCertificatesClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for App Service Certificate creation.")

	resGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)
	location := d.Get("location").(string)
	pfxString := d.Get("base_64_encoded_pfx_file").(string)
	password := d.Get("pfx_password").(string)

	pfxBlob, err := base64.StdEncoding.DecodeString(pfxString)

	certificate := web.Certificate{
		Name:     &name,
		Location: &location,
		CertificateProperties: &web.CertificateProperties{
			PfxBlob:  &pfxBlob,
			Password: &password,
		},
	}

	certificateResult, err := client.CreateOrUpdate(ctx, resGroup, name, certificate)
	if err != nil {
		return err
	}

	d.Set("thumbprint", certificateResult.CertificateProperties.Thumbprint)
	d.SetId(*certificateResult.ID)
	return nil
}

func resourceAppServiceCertificateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).appServicesCertificatesClient
	ctx := meta.(*ArmClient).StopContext

	resGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return err
	}

	d.Set("thumbprint", read.CertificateProperties.Thumbprint)
	d.SetId(*read.ID)
	return nil
}

func resourceAppServiceCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
