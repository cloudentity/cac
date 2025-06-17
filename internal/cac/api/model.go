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

type Patch interface {
	GetData() models.Rfc7396PatchOperation
	GetExtensions() any
	Merge(other Patch) error
}

type PatchImpl[T any] struct {
	Data models.Rfc7396PatchOperation `json:"data,omitempty"`
	Ext  *T                           `json:"ext,omitempty"`
}

type ServerPatch PatchImpl[ServerExtensions]

var _ Patch = &ServerPatch{}

func (sp *ServerPatch) GetData() models.Rfc7396PatchOperation {
	return sp.Data
}
func (tp *ServerPatch) GetExtensions() any {
	return tp.Ext
}
func (sp *ServerPatch) Merge(other Patch) error {
	if err := mergo.Merge(&sp.Data, other.GetData(), mergo.WithOverride); err != nil {
		return err
	}

	if err := mergo.Merge(sp.Ext, other.GetExtensions(), mergo.WithOverride); err != nil {
		return err
	}

	return nil
}

type TenantPatch PatchImpl[TenantExtensions]

func (tp *TenantPatch) GetData() models.Rfc7396PatchOperation {
	return tp.Data
}
func (tp *TenantPatch) GetExtensions() any {
	return tp.Ext
}
func (sp *TenantPatch) Merge(other Patch) error {
	if err := mergo.Merge(&sp.Data, other.GetData(), mergo.WithOverride); err != nil {
		return err
	}

	if err := mergo.Merge(sp.Ext, other.GetExtensions(), mergo.WithOverride); err != nil {
		return err
	}

	return nil
}
