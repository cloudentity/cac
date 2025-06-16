package api

import (
	"github.com/cloudentity/acp-client-go/clients/hub/models"
	smodels "github.com/cloudentity/acp-client-go/clients/system/models"
	"github.com/imdario/mergo"
)

type ServerExtensions struct {
	Secrets map[string]*smodels.Secret `json:"secrets,omitempty"`
}

type TenantExtensions struct {
	Servers map[string]ServerExtensions `json:"servers,omitempty"`
}

func (te *TenantExtensions) GetServerExtensions(serverID string) *ServerExtensions {
	if te.Servers == nil {
		return nil
	}

	if ext, ok := te.Servers[serverID]; ok {
		return &ext
	}

	return nil
}

type Patch[T any] struct {
	Data models.Rfc7396PatchOperation `json:"data,omitempty"`
	Ext  *T                           `json:"ext,omitempty"`
}

type PatchInterface interface {
	GetData() models.Rfc7396PatchOperation
	GetExtensions() any
	Merge(other PatchInterface) error
}

type ServerPatch Patch[ServerExtensions]

func (sp *ServerPatch) GetData() models.Rfc7396PatchOperation {
	return sp.Data
}
func (tp *ServerPatch) GetExtensions() any {
	return tp.Ext
}
func (sp *ServerPatch) Merge(other PatchInterface) error {
	if err := mergo.Merge(&sp.Data, other.GetData(), mergo.WithOverride); err != nil {
		return err
	}

	if err := mergo.Merge(sp.Ext, other.GetExtensions(), mergo.WithOverride); err != nil {
		return err
	}

	return nil
}

type TenantPatch Patch[TenantExtensions]

func (tp *TenantPatch) GetData() models.Rfc7396PatchOperation {
	return tp.Data
}
func (tp *TenantPatch) GetExtensions() any {
	return tp.Ext
}
func (sp *TenantPatch) Merge(other PatchInterface) error {
	if err := mergo.Merge(&sp.Data, other.GetData(), mergo.WithOverride); err != nil {
		return err
	}

	if err := mergo.Merge(sp.Ext, other.GetExtensions(), mergo.WithOverride); err != nil {
		return err
	}

	return nil
}