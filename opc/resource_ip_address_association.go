package opc

import (
	"fmt"

	"github.com/hashicorp/go-oracle-terraform/client"
	"github.com/hashicorp/go-oracle-terraform/compute"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceOPCIPAddressAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceOPCIPAddressAssociationCreate,
		Read:   resourceOPCIPAddressAssociationRead,
		Delete: resourceOPCIPAddressAssociationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ip_address_reservation": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vnic": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"tags": tagsForceNewSchema(),
			"uri": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOPCIPAddressAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	computeClient, err := meta.(*Client).getComputeClient()
	if err != nil {
		return err
	}
	resClient := computeClient.IPAddressAssociations()

	input := compute.CreateIPAddressAssociationInput{
		Name: d.Get("name").(string),
	}

	if ipAddressReservation, ok := d.GetOk("ip_address_reservation"); ok {
		input.IPAddressReservation = ipAddressReservation.(string)
	}

	if vnic, ok := d.GetOk("vnic"); ok {
		input.Vnic = vnic.(string)
	}

	tags := getStringList(d, "tags")
	if len(tags) != 0 {
		input.Tags = tags
	}

	if description, ok := d.GetOk("description"); ok {
		input.Description = description.(string)
	}

	info, err := resClient.CreateIPAddressAssociation(&input)
	if err != nil {
		return fmt.Errorf("Error creating IP Address Association: %s", err)
	}

	d.SetId(info.Name)
	return resourceOPCIPAddressAssociationRead(d, meta)
}

func resourceOPCIPAddressAssociationRead(d *schema.ResourceData, meta interface{}) error {
	computeClient, err := meta.(*Client).getComputeClient()
	if err != nil {
		return err
	}
	resClient := computeClient.IPAddressAssociations()
	name := d.Id()

	getInput := compute.GetIPAddressAssociationInput{
		Name: name,
	}
	result, err := resClient.GetIPAddressAssociation(&getInput)
	if err != nil {
		// IP Address Association does not exist
		if client.WasNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading IP Address Association %s: %s", name, err)
	}
	if result == nil {
		d.SetId("")
		return fmt.Errorf("Error reading IP Address Association %s: %s", name, err)
	}

	d.Set("name", result.Name)
	d.Set("ip_address_reservation", result.IPAddressReservation)
	d.Set("vnic", result.Vnic)
	d.Set("description", result.Description)
	d.Set("uri", result.URI)
	if err := setStringList(d, "tags", result.Tags); err != nil {
		return err
	}
	return nil
}

func resourceOPCIPAddressAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	computeClient, err := meta.(*Client).getComputeClient()
	if err != nil {
		return err
	}
	resClient := computeClient.IPAddressAssociations()
	name := d.Id()

	input := compute.DeleteIPAddressAssociationInput{
		Name: name,
	}
	if err := resClient.DeleteIPAddressAssociation(&input); err != nil {
		return fmt.Errorf("Error deleting IP Address Association: %s", err)
	}
	return nil
}
