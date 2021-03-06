// Copyright (c) 2017, 2020, Oracle and/or its affiliates. All rights reserved.
// Licensed under the Mozilla Public License v2.0

package oci

import (
	"context"
	"encoding/base64"
	"io/ioutil"

	"github.com/hashicorp/terraform/helper/schema"
	oci_database "github.com/oracle/oci-go-sdk/database"
)

func init() {
	RegisterDatasource("oci_database_autonomous_data_warehouse_wallet", DatabaseAutonomousDataWarehouseWalletDataSource())
}

func DatabaseAutonomousDataWarehouseWalletDataSource() *schema.Resource {
	return &schema.Resource{
		Read: readSingularDatabaseAutonomousDataWarehouseWallet,
		Schema: map[string]*schema.Schema{
			"autonomous_data_warehouse_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"base64_encode_content": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			// Computed
			"content": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func readSingularDatabaseAutonomousDataWarehouseWallet(d *schema.ResourceData, m interface{}) error {
	sync := &DatabaseAutonomousDataWarehouseWalletDataSourceCrud{}
	sync.D = d
	sync.Client = m.(*OracleClients).databaseClient()

	return ReadResource(sync)
}

type DatabaseAutonomousDataWarehouseWalletDataSourceCrud struct {
	D      *schema.ResourceData
	Client *oci_database.DatabaseClient
	Res    *[]byte
}

func (s *DatabaseAutonomousDataWarehouseWalletDataSourceCrud) VoidState() {
	s.D.SetId("")
}

func (s *DatabaseAutonomousDataWarehouseWalletDataSourceCrud) Get() error {
	request := oci_database.GenerateAutonomousDataWarehouseWalletRequest{}

	if autonomousDataWarehouseId, ok := s.D.GetOkExists("autonomous_data_warehouse_id"); ok {
		tmp := autonomousDataWarehouseId.(string)
		request.AutonomousDataWarehouseId = &tmp
	}

	if password, ok := s.D.GetOkExists("password"); ok {
		tmp := password.(string)
		request.Password = &tmp
	}

	request.RequestMetadata.RetryPolicy = getRetryPolicy(false, "database")

	response, err := s.Client.GenerateAutonomousDataWarehouseWallet(context.Background(), request)
	if err != nil {
		return err
	}

	if response.Content != nil {
		defer response.Content.Close()
		if contentBytes, err := ioutil.ReadAll(response.Content); err == nil {
			s.Res = &contentBytes
		} else {
			return err
		}
	}

	return nil
}

func (s *DatabaseAutonomousDataWarehouseWalletDataSourceCrud) SetData() error {
	if s.Res == nil {
		return nil
	}

	s.D.SetId(GenerateDataSourceID())

	base64EncodeContent := false
	if tmp, ok := s.D.GetOkExists("base64_encode_content"); ok {
		base64EncodeContent = tmp.(bool)
	}

	if base64EncodeContent {
		s.D.Set("content", base64.StdEncoding.EncodeToString(*s.Res))
	} else {
		s.D.Set("content", string(*s.Res))
	}

	return nil
}
