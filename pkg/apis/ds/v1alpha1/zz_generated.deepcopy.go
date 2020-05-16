// +build !ignore_autogenerated

// Code generated by operator-sdk-v0.17.0-x86_64-linux-gnu. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Xdaemonset) DeepCopyInto(out *Xdaemonset) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Xdaemonset.
func (in *Xdaemonset) DeepCopy() *Xdaemonset {
	if in == nil {
		return nil
	}
	out := new(Xdaemonset)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Xdaemonset) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *XdaemonsetList) DeepCopyInto(out *XdaemonsetList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Xdaemonset, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new XdaemonsetList.
func (in *XdaemonsetList) DeepCopy() *XdaemonsetList {
	if in == nil {
		return nil
	}
	out := new(XdaemonsetList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *XdaemonsetList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *XdaemonsetSpec) DeepCopyInto(out *XdaemonsetSpec) {
	*out = *in
	in.DaemonSetSpec.DeepCopyInto(&out.DaemonSetSpec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new XdaemonsetSpec.
func (in *XdaemonsetSpec) DeepCopy() *XdaemonsetSpec {
	if in == nil {
		return nil
	}
	out := new(XdaemonsetSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *XdaemonsetStatus) DeepCopyInto(out *XdaemonsetStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new XdaemonsetStatus.
func (in *XdaemonsetStatus) DeepCopy() *XdaemonsetStatus {
	if in == nil {
		return nil
	}
	out := new(XdaemonsetStatus)
	in.DeepCopyInto(out)
	return out
}
