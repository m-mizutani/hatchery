// Code generated from Pkl module `org.github.m_mizutani.hatchery.config`. DO NOT EDIT.
package config

type Action interface {
	GetId() string

	GetTags() *[]string

	GetBucket() string

	GetPrefix() *string
}
