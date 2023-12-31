// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
	v1alpha1 "open-cluster-management.io/api/addon/v1alpha1"
	scheme "open-cluster-management.io/api/client/addon/clientset/versioned/scheme"
)

// AddOnTemplatesGetter has a method to return a AddOnTemplateInterface.
// A group's client should implement this interface.
type AddOnTemplatesGetter interface {
	AddOnTemplates() AddOnTemplateInterface
}

// AddOnTemplateInterface has methods to work with AddOnTemplate resources.
type AddOnTemplateInterface interface {
	Create(ctx context.Context, addOnTemplate *v1alpha1.AddOnTemplate, opts v1.CreateOptions) (*v1alpha1.AddOnTemplate, error)
	Update(ctx context.Context, addOnTemplate *v1alpha1.AddOnTemplate, opts v1.UpdateOptions) (*v1alpha1.AddOnTemplate, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.AddOnTemplate, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.AddOnTemplateList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.AddOnTemplate, err error)
	AddOnTemplateExpansion
}

// addOnTemplates implements AddOnTemplateInterface
type addOnTemplates struct {
	client rest.Interface
}

// newAddOnTemplates returns a AddOnTemplates
func newAddOnTemplates(c *AddonV1alpha1Client) *addOnTemplates {
	return &addOnTemplates{
		client: c.RESTClient(),
	}
}

// Get takes name of the addOnTemplate, and returns the corresponding addOnTemplate object, and an error if there is any.
func (c *addOnTemplates) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.AddOnTemplate, err error) {
	result = &v1alpha1.AddOnTemplate{}
	err = c.client.Get().
		Resource("addontemplates").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of AddOnTemplates that match those selectors.
func (c *addOnTemplates) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.AddOnTemplateList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.AddOnTemplateList{}
	err = c.client.Get().
		Resource("addontemplates").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested addOnTemplates.
func (c *addOnTemplates) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("addontemplates").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a addOnTemplate and creates it.  Returns the server's representation of the addOnTemplate, and an error, if there is any.
func (c *addOnTemplates) Create(ctx context.Context, addOnTemplate *v1alpha1.AddOnTemplate, opts v1.CreateOptions) (result *v1alpha1.AddOnTemplate, err error) {
	result = &v1alpha1.AddOnTemplate{}
	err = c.client.Post().
		Resource("addontemplates").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(addOnTemplate).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a addOnTemplate and updates it. Returns the server's representation of the addOnTemplate, and an error, if there is any.
func (c *addOnTemplates) Update(ctx context.Context, addOnTemplate *v1alpha1.AddOnTemplate, opts v1.UpdateOptions) (result *v1alpha1.AddOnTemplate, err error) {
	result = &v1alpha1.AddOnTemplate{}
	err = c.client.Put().
		Resource("addontemplates").
		Name(addOnTemplate.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(addOnTemplate).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the addOnTemplate and deletes it. Returns an error if one occurs.
func (c *addOnTemplates) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("addontemplates").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *addOnTemplates) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("addontemplates").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched addOnTemplate.
func (c *addOnTemplates) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.AddOnTemplate, err error) {
	result = &v1alpha1.AddOnTemplate{}
	err = c.client.Patch(pt).
		Resource("addontemplates").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
