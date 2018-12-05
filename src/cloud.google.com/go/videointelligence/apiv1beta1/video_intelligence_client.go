// Copyright 2017, Google Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// AUTO-GENERATED CODE. DO NOT EDIT.

package videointelligence

import (
	"time"

	"cloud.google.com/go/internal/version"
	"cloud.google.com/go/longrunning"
	gax "github.com/googleapis/gax-go"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"google.golang.org/api/transport"
	videointelligencepb "google.golang.org/genproto/googleapis/cloud/videointelligence/v1beta1"
	longrunningpb "google.golang.org/genproto/googleapis/longrunning"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// VideoIntelligenceCallOptions contains the retry settings for each method of VideoIntelligenceClient.
type VideoIntelligenceCallOptions struct {
	AnnotateVideo []gax.CallOption
}

func defaultVideoIntelligenceClientOptions() []option.ClientOption {
	return []option.ClientOption{
		option.WithEndpoint("videointelligence.googleapis.com:443"),
		option.WithScopes(DefaultAuthScopes()...),
	}
}

func defaultVideoIntelligenceCallOptions() *VideoIntelligenceCallOptions {
	retry := map[[2]string][]gax.CallOption{
		{"default", "idempotent"}: {
			gax.WithRetry(func() gax.Retryer {
				return gax.OnCodes([]codes.Code{
					codes.DeadlineExceeded,
					codes.Unavailable,
				}, gax.Backoff{
					Initial:    1000 * time.Millisecond,
					Max:        120000 * time.Millisecond,
					Multiplier: 2.5,
				})
			}),
		},
	}
	return &VideoIntelligenceCallOptions{
		AnnotateVideo: retry[[2]string{"default", "idempotent"}],
	}
}

// VideoIntelligenceClient is a client for interacting with Google Cloud Video Intelligence API.
type VideoIntelligenceClient struct {
	// The connection to the service.
	conn *grpc.ClientConn

	// The gRPC API client.
	videoIntelligenceClient videointelligencepb.VideoIntelligenceServiceClient

	// The call options for this service.
	CallOptions *VideoIntelligenceCallOptions

	// The metadata to be sent with each request.
	xGoogHeader []string
}

// NewVideoIntelligenceClient creates a new video intelligence service client.
//
// Service that implements Google Cloud Video Intelligence API.
func NewVideoIntelligenceClient(ctx context.Context, opts ...option.ClientOption) (*VideoIntelligenceClient, error) {
	conn, err := transport.DialGRPC(ctx, append(defaultVideoIntelligenceClientOptions(), opts...)...)
	if err != nil {
		return nil, err
	}
	c := &VideoIntelligenceClient{
		conn:        conn,
		CallOptions: defaultVideoIntelligenceCallOptions(),

		videoIntelligenceClient: videointelligencepb.NewVideoIntelligenceServiceClient(conn),
	}
	c.SetGoogleClientInfo()
	return c, nil
}

// Connection returns the client's connection to the API service.
func (c *VideoIntelligenceClient) Connection() *grpc.ClientConn {
	return c.conn
}

// Close closes the connection to the API service. The user should invoke this when
// the client is no longer required.
func (c *VideoIntelligenceClient) Close() error {
	return c.conn.Close()
}

// SetGoogleClientInfo sets the name and version of the application in
// the `x-goog-api-client` header passed on each request. Intended for
// use by Google-written clients.
func (c *VideoIntelligenceClient) SetGoogleClientInfo(keyval ...string) {
	kv := append([]string{"gl-go", version.Go()}, keyval...)
	kv = append(kv, "gapic", version.Repo, "gax", gax.Version, "grpc", grpc.Version)
	c.xGoogHeader = []string{gax.XGoogHeader(kv...)}
}

// AnnotateVideo performs asynchronous video annotation. Progress and results can be
// retrieved through the `google.longrunning.Operations` interface.
// `Operation.metadata` contains `AnnotateVideoProgress` (progress).
// `Operation.response` contains `AnnotateVideoResponse` (results).
func (c *VideoIntelligenceClient) AnnotateVideo(ctx context.Context, req *videointelligencepb.AnnotateVideoRequest, opts ...gax.CallOption) (*AnnotateVideoOperation, error) {
	ctx = insertXGoog(ctx, c.xGoogHeader)
	opts = append(c.CallOptions.AnnotateVideo[0:len(c.CallOptions.AnnotateVideo):len(c.CallOptions.AnnotateVideo)], opts...)
	var resp *longrunningpb.Operation
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.videoIntelligenceClient.AnnotateVideo(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return &AnnotateVideoOperation{
		lro:         longrunning.InternalNewOperation(c.Connection(), resp),
		xGoogHeader: c.xGoogHeader,
	}, nil
}

// AnnotateVideoOperation manages a long-running operation from AnnotateVideo.
type AnnotateVideoOperation struct {
	lro *longrunning.Operation

	// The metadata to be sent with each request.
	xGoogHeader []string
}

// AnnotateVideoOperation returns a new AnnotateVideoOperation from a given name.
// The name must be that of a previously created AnnotateVideoOperation, possibly from a different process.
func (c *VideoIntelligenceClient) AnnotateVideoOperation(name string) *AnnotateVideoOperation {
	return &AnnotateVideoOperation{
		lro:         longrunning.InternalNewOperation(c.Connection(), &longrunningpb.Operation{Name: name}),
		xGoogHeader: c.xGoogHeader,
	}
}

// Wait blocks until the long-running operation is completed, returning the response and any errors encountered.
//
// See documentation of Poll for error-handling information.
func (op *AnnotateVideoOperation) Wait(ctx context.Context) (*videointelligencepb.AnnotateVideoResponse, error) {
	var resp videointelligencepb.AnnotateVideoResponse
	ctx = insertXGoog(ctx, op.xGoogHeader)
	if err := op.lro.Wait(ctx, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Poll fetches the latest state of the long-running operation.
//
// Poll also fetches the latest metadata, which can be retrieved by Metadata.
//
// If Poll fails, the error is returned and op is unmodified. If Poll succeeds and
// the operation has completed with failure, the error is returned and op.Done will return true.
// If Poll succeeds and the operation has completed successfully,
// op.Done will return true, and the response of the operation is returned.
// If Poll succeeds and the operation has not completed, the returned response and error are both nil.
func (op *AnnotateVideoOperation) Poll(ctx context.Context) (*videointelligencepb.AnnotateVideoResponse, error) {
	var resp videointelligencepb.AnnotateVideoResponse
	ctx = insertXGoog(ctx, op.xGoogHeader)
	if err := op.lro.Poll(ctx, &resp); err != nil {
		return nil, err
	}
	if !op.Done() {
		return nil, nil
	}
	return &resp, nil
}

// Metadata returns metadata associated with the long-running operation.
// Metadata itself does not contact the server, but Poll does.
// To get the latest metadata, call this method after a successful call to Poll.
// If the metadata is not available, the returned metadata and error are both nil.
func (op *AnnotateVideoOperation) Metadata() (*videointelligencepb.AnnotateVideoProgress, error) {
	var meta videointelligencepb.AnnotateVideoProgress
	if err := op.lro.Metadata(&meta); err == longrunning.ErrNoMetadata {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &meta, nil
}

// Done reports whether the long-running operation has completed.
func (op *AnnotateVideoOperation) Done() bool {
	return op.lro.Done()
}

// Name returns the name of the long-running operation.
// The name is assigned by the server and is unique within the service from which the operation is created.
func (op *AnnotateVideoOperation) Name() string {
	return op.lro.Name()
}
