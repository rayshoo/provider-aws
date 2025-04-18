/*
Copyright 2021 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by ack-generate. DO NOT EDIT.

package methodresponse

import (
	"context"

	svcapi "github.com/aws/aws-sdk-go/service/apigateway"
	svcsdk "github.com/aws/aws-sdk-go/service/apigateway"
	svcsdkapi "github.com/aws/aws-sdk-go/service/apigateway/apigatewayiface"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/meta"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	cpresource "github.com/crossplane/crossplane-runtime/pkg/resource"

	svcapitypes "github.com/crossplane-contrib/provider-aws/apis/apigateway/v1alpha1"
	connectaws "github.com/crossplane-contrib/provider-aws/pkg/utils/connect/aws"
	errorutils "github.com/crossplane-contrib/provider-aws/pkg/utils/errors"
)

const (
	errUnexpectedObject = "managed resource is not an MethodResponse resource"

	errCreateSession = "cannot create a new session"
	errCreate        = "cannot create MethodResponse in AWS"
	errUpdate        = "cannot update MethodResponse in AWS"
	errDescribe      = "failed to describe MethodResponse"
	errDelete        = "failed to delete MethodResponse"
)

type connector struct {
	kube client.Client
	opts []option
}

func (c *connector) Connect(ctx context.Context, cr *svcapitypes.MethodResponse) (managed.TypedExternalClient[*svcapitypes.MethodResponse], error) {
	sess, err := connectaws.GetConfigV1(ctx, c.kube, cr, cr.Spec.ForProvider.Region)
	if err != nil {
		return nil, errors.Wrap(err, errCreateSession)
	}
	return newExternal(c.kube, svcapi.New(sess), c.opts), nil
}

func (e *external) Observe(ctx context.Context, cr *svcapitypes.MethodResponse) (managed.ExternalObservation, error) {
	if meta.GetExternalName(cr) == "" {
		return managed.ExternalObservation{
			ResourceExists: false,
		}, nil
	}
	input := GenerateGetMethodResponseInput(cr)
	if err := e.preObserve(ctx, cr, input); err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, "pre-observe failed")
	}
	resp, err := e.client.GetMethodResponseWithContext(ctx, input)
	if err != nil {
		return managed.ExternalObservation{ResourceExists: false}, errorutils.Wrap(cpresource.Ignore(IsNotFound, err), errDescribe)
	}
	currentSpec := cr.Spec.ForProvider.DeepCopy()
	if err := e.lateInitialize(&cr.Spec.ForProvider, resp); err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, "late-init failed")
	}
	GenerateMethodResponse(resp).Status.AtProvider.DeepCopyInto(&cr.Status.AtProvider)
	upToDate := true
	diff := ""
	if !meta.WasDeleted(cr) { // There is no need to run isUpToDate if the resource is deleted
		upToDate, diff, err = e.isUpToDate(ctx, cr, resp)
		if err != nil {
			return managed.ExternalObservation{}, errors.Wrap(err, "isUpToDate check failed")
		}
	}
	return e.postObserve(ctx, cr, resp, managed.ExternalObservation{
		ResourceExists:          true,
		ResourceUpToDate:        upToDate,
		Diff:                    diff,
		ResourceLateInitialized: !cmp.Equal(&cr.Spec.ForProvider, currentSpec),
	}, nil)
}

func (e *external) Create(ctx context.Context, cr *svcapitypes.MethodResponse) (managed.ExternalCreation, error) {
	cr.Status.SetConditions(xpv1.Creating())
	input := GeneratePutMethodResponseInput(cr)
	if err := e.preCreate(ctx, cr, input); err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, "pre-create failed")
	}
	resp, err := e.client.PutMethodResponseWithContext(ctx, input)
	if err != nil {
		return managed.ExternalCreation{}, errorutils.Wrap(err, errCreate)
	}

	if resp.ResponseModels != nil {
		f0 := map[string]*string{}
		for f0key, f0valiter := range resp.ResponseModels {
			var f0val string
			f0val = *f0valiter
			f0[f0key] = &f0val
		}
		cr.Spec.ForProvider.ResponseModels = f0
	} else {
		cr.Spec.ForProvider.ResponseModels = nil
	}
	if resp.ResponseParameters != nil {
		f1 := map[string]*bool{}
		for f1key, f1valiter := range resp.ResponseParameters {
			var f1val bool
			f1val = *f1valiter
			f1[f1key] = &f1val
		}
		cr.Spec.ForProvider.ResponseParameters = f1
	} else {
		cr.Spec.ForProvider.ResponseParameters = nil
	}
	if resp.StatusCode != nil {
		cr.Spec.ForProvider.StatusCode = resp.StatusCode
	} else {
		cr.Spec.ForProvider.StatusCode = nil
	}

	return e.postCreate(ctx, cr, resp, managed.ExternalCreation{}, err)
}

func (e *external) Update(ctx context.Context, cr *svcapitypes.MethodResponse) (managed.ExternalUpdate, error) {
	input := GenerateUpdateMethodResponseInput(cr)
	if err := e.preUpdate(ctx, cr, input); err != nil {
		return managed.ExternalUpdate{}, errors.Wrap(err, "pre-update failed")
	}
	resp, err := e.client.UpdateMethodResponseWithContext(ctx, input)
	return e.postUpdate(ctx, cr, resp, managed.ExternalUpdate{}, errorutils.Wrap(err, errUpdate))
}

func (e *external) Delete(ctx context.Context, cr *svcapitypes.MethodResponse) (managed.ExternalDelete, error) {
	cr.Status.SetConditions(xpv1.Deleting())
	input := GenerateDeleteMethodResponseInput(cr)
	ignore, err := e.preDelete(ctx, cr, input)
	if err != nil {
		return managed.ExternalDelete{}, errors.Wrap(err, "pre-delete failed")
	}
	if ignore {
		return managed.ExternalDelete{}, nil
	}
	resp, err := e.client.DeleteMethodResponseWithContext(ctx, input)
	return e.postDelete(ctx, cr, resp, errorutils.Wrap(cpresource.Ignore(IsNotFound, err), errDelete))
}

func (e *external) Disconnect(ctx context.Context) error {
	// Unimplemented, required by newer versions of crossplane-runtime
	return nil
}

type option func(*external)

func newExternal(kube client.Client, client svcsdkapi.APIGatewayAPI, opts []option) *external {
	e := &external{
		kube:           kube,
		client:         client,
		preObserve:     nopPreObserve,
		postObserve:    nopPostObserve,
		lateInitialize: nopLateInitialize,
		isUpToDate:     alwaysUpToDate,
		preCreate:      nopPreCreate,
		postCreate:     nopPostCreate,
		preDelete:      nopPreDelete,
		postDelete:     nopPostDelete,
		preUpdate:      nopPreUpdate,
		postUpdate:     nopPostUpdate,
	}
	for _, f := range opts {
		f(e)
	}
	return e
}

type external struct {
	kube           client.Client
	client         svcsdkapi.APIGatewayAPI
	preObserve     func(context.Context, *svcapitypes.MethodResponse, *svcsdk.GetMethodResponseInput) error
	postObserve    func(context.Context, *svcapitypes.MethodResponse, *svcsdk.MethodResponse, managed.ExternalObservation, error) (managed.ExternalObservation, error)
	lateInitialize func(*svcapitypes.MethodResponseParameters, *svcsdk.MethodResponse) error
	isUpToDate     func(context.Context, *svcapitypes.MethodResponse, *svcsdk.MethodResponse) (bool, string, error)
	preCreate      func(context.Context, *svcapitypes.MethodResponse, *svcsdk.PutMethodResponseInput) error
	postCreate     func(context.Context, *svcapitypes.MethodResponse, *svcsdk.MethodResponse, managed.ExternalCreation, error) (managed.ExternalCreation, error)
	preDelete      func(context.Context, *svcapitypes.MethodResponse, *svcsdk.DeleteMethodResponseInput) (bool, error)
	postDelete     func(context.Context, *svcapitypes.MethodResponse, *svcsdk.DeleteMethodResponseOutput, error) (managed.ExternalDelete, error)
	preUpdate      func(context.Context, *svcapitypes.MethodResponse, *svcsdk.UpdateMethodResponseInput) error
	postUpdate     func(context.Context, *svcapitypes.MethodResponse, *svcsdk.MethodResponse, managed.ExternalUpdate, error) (managed.ExternalUpdate, error)
}

func nopPreObserve(context.Context, *svcapitypes.MethodResponse, *svcsdk.GetMethodResponseInput) error {
	return nil
}

func nopPostObserve(_ context.Context, _ *svcapitypes.MethodResponse, _ *svcsdk.MethodResponse, obs managed.ExternalObservation, err error) (managed.ExternalObservation, error) {
	return obs, err
}
func nopLateInitialize(*svcapitypes.MethodResponseParameters, *svcsdk.MethodResponse) error {
	return nil
}
func alwaysUpToDate(context.Context, *svcapitypes.MethodResponse, *svcsdk.MethodResponse) (bool, string, error) {
	return true, "", nil
}

func nopPreCreate(context.Context, *svcapitypes.MethodResponse, *svcsdk.PutMethodResponseInput) error {
	return nil
}
func nopPostCreate(_ context.Context, _ *svcapitypes.MethodResponse, _ *svcsdk.MethodResponse, cre managed.ExternalCreation, err error) (managed.ExternalCreation, error) {
	return cre, err
}
func nopPreDelete(context.Context, *svcapitypes.MethodResponse, *svcsdk.DeleteMethodResponseInput) (bool, error) {
	return false, nil
}
func nopPostDelete(_ context.Context, _ *svcapitypes.MethodResponse, _ *svcsdk.DeleteMethodResponseOutput, err error) (managed.ExternalDelete, error) {
	return managed.ExternalDelete{}, err
}
func nopPreUpdate(context.Context, *svcapitypes.MethodResponse, *svcsdk.UpdateMethodResponseInput) error {
	return nil
}
func nopPostUpdate(_ context.Context, _ *svcapitypes.MethodResponse, _ *svcsdk.MethodResponse, upd managed.ExternalUpdate, err error) (managed.ExternalUpdate, error) {
	return upd, err
}
