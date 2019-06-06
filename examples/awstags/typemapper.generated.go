// Code generated by "typemapper "; DO NOT EDIT.

// +build !typemapper

package awstags

import (
	datasync "github.com/aws/aws-sdk-go/service/datasync"
	ec2 "github.com/aws/aws-sdk-go/service/ec2"
	elbv2 "github.com/aws/aws-sdk-go/service/elbv2"
)

func ELBv2TagToEC2Tag(src *elbv2.Tag, dst *ec2.Tag) error {
	if dst == nil {
		return nil
	}
	dst.Key = src.Key
	dst.Value = src.Value
	return nil
}
func (src *myTag) DataSyncTag(dst *datasync.TagListEntry) error {
	if dst == nil {
		return nil
	}
	dst.Key = &src.Key
	dst.Value = &src.Value
	return nil
}
func (src *myTag) EC2Tag() {
	if dst == nil {
		return
	}
	dst.Key = &src.Key
	dst.Value = &src.Value
	return
}
func EC2TagToDataSyncTag(src *ec2.Tag, dst *datasync.TagListEntry) error {
	if dst == nil {
		return nil
	}
	dst.Key = src.Key
	dst.Value = src.Value
	return nil
}