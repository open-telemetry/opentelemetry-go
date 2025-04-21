// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Code generated from semantic convention specification. DO NOT EDIT.

package semconv // import "go.opentelemetry.io/otel/semconv/v1.31.0/vcs"

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// ChangeStateAttr is an attribute conforming to the vcs.change.state semantic
// conventions. It represents the state of the change (pull request/merge
// request/changelist).
type ChangeStateAttr string

var (
	// ChangeStateOpen is the open means the change is currently active and under
	// review. It hasn't been merged into the target branch yet, and it's still
	// possible to make changes or add comments.
	ChangeStateOpen ChangeStateAttr = "open"
	// ChangeStateWip is the WIP (work-in-progress, draft) means the change is still
	// in progress and not yet ready for a full review. It might still undergo
	// significant changes.
	ChangeStateWip ChangeStateAttr = "wip"
	// ChangeStateClosed is the closed means the merge request has been closed
	// without merging. This can happen for various reasons, such as the changes
	// being deemed unnecessary, the issue being resolved in another way, or the
	// author deciding to withdraw the request.
	ChangeStateClosed ChangeStateAttr = "closed"
	// ChangeStateMerged is the merged indicates that the change has been
	// successfully integrated into the target codebase.
	ChangeStateMerged ChangeStateAttr = "merged"
)

// LineChangeTypeAttr is an attribute conforming to the vcs.line_change.type
// semantic conventions. It represents the type of line change being measured on
// a branch or change.
type LineChangeTypeAttr string

var (
	// LineChangeTypeAdded is the how many lines were added.
	LineChangeTypeAdded LineChangeTypeAttr = "added"
	// LineChangeTypeRemoved is the how many lines were removed.
	LineChangeTypeRemoved LineChangeTypeAttr = "removed"
)

// RefBaseTypeAttr is an attribute conforming to the vcs.ref.base.type semantic
// conventions. It represents the type of the [reference] in the repository.
//
// [reference]: https://git-scm.com/docs/gitglossary#def_ref
type RefBaseTypeAttr string

var (
	// RefBaseTypeBranch is the [branch].
	//
	// [branch]: https://git-scm.com/docs/gitglossary#Documentation/gitglossary.txt-aiddefbranchabranch
	RefBaseTypeBranch RefBaseTypeAttr = "branch"
	// RefBaseTypeTag is the [tag].
	//
	// [tag]: https://git-scm.com/docs/gitglossary#Documentation/gitglossary.txt-aiddeftagatag
	RefBaseTypeTag RefBaseTypeAttr = "tag"
)

// RefHeadTypeAttr is an attribute conforming to the vcs.ref.head.type semantic
// conventions. It represents the type of the [reference] in the repository.
//
// [reference]: https://git-scm.com/docs/gitglossary#def_ref
type RefHeadTypeAttr string

var (
	// RefHeadTypeBranch is the [branch].
	//
	// [branch]: https://git-scm.com/docs/gitglossary#Documentation/gitglossary.txt-aiddefbranchabranch
	RefHeadTypeBranch RefHeadTypeAttr = "branch"
	// RefHeadTypeTag is the [tag].
	//
	// [tag]: https://git-scm.com/docs/gitglossary#Documentation/gitglossary.txt-aiddeftagatag
	RefHeadTypeTag RefHeadTypeAttr = "tag"
)

// RefTypeAttr is an attribute conforming to the vcs.ref.type semantic
// conventions. It represents the type of the [reference] in the repository.
//
// [reference]: https://git-scm.com/docs/gitglossary#def_ref
type RefTypeAttr string

var (
	// RefTypeBranch is the [branch].
	//
	// [branch]: https://git-scm.com/docs/gitglossary#Documentation/gitglossary.txt-aiddefbranchabranch
	RefTypeBranch RefTypeAttr = "branch"
	// RefTypeTag is the [tag].
	//
	// [tag]: https://git-scm.com/docs/gitglossary#Documentation/gitglossary.txt-aiddeftagatag
	RefTypeTag RefTypeAttr = "tag"
)

// RevisionDeltaDirectionAttr is an attribute conforming to the
// vcs.revision_delta.direction semantic conventions. It represents the type of
// revision comparison.
type RevisionDeltaDirectionAttr string

var (
	// RevisionDeltaDirectionBehind is the how many revisions the change is behind
	// the target ref.
	RevisionDeltaDirectionBehind RevisionDeltaDirectionAttr = "behind"
	// RevisionDeltaDirectionAhead is the how many revisions the change is ahead of
	// the target ref.
	RevisionDeltaDirectionAhead RevisionDeltaDirectionAttr = "ahead"
)

// VCSChangeCount is an instrument used to record metric values conforming to the
// "vcs.change.count" semantic conventions. It represents the number of changes
// (pull requests/merge requests/changelists) in a repository, categorized by
// their state (e.g. open or merged).
type ChangeCount struct {
	inst metric.Int64UpDownCounter
}

// NewChangeCount returns a new ChangeCount instrument.
func NewChangeCount(m metric.Meter) (ChangeCount, error) {
	i, err := m.Int64UpDownCounter(
	    "vcs.change.count",
	    metric.WithDescription("The number of changes (pull requests/merge requests/changelists) in a repository, categorized by their state (e.g. open or merged)"),
	    metric.WithUnit("{change}"),
	)
	if err != nil {
	    return ChangeCount{}, err
	}
	return ChangeCount{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (ChangeCount) Name() string {
	return "vcs.change.count"
}

// Unit returns the semantic convention unit of the instrument
func (ChangeCount) Unit() string {
	return "{change}"
}

// Description returns the semantic convention description of the instrument
func (ChangeCount) Description() string {
	return "The number of changes (pull requests/merge requests/changelists) in a repository, categorized by their state (e.g. open or merged)"
}

// Add adds incr to the existing count.
//
// The vcsChangeState is the the state of the change (pull request/merge
// request/changelist).
//
// The vcsRepositoryUrlFull is the the [canonical URL] of the repository
// providing the complete HTTP(S) address in order to locate and identify the
// repository through a browser.
//
// All additional attrs passed are included in the recorded value.
//
// [canonical URL]: https://support.google.com/webmasters/answer/10347851?hl=en#:~:text=A%20canonical%20URL%20is%20the,Google%20chooses%20one%20as%20canonical.
func (m ChangeCount) Add(
    ctx context.Context,
    incr int64,
	changeState ChangeStateAttr,
	repositoryUrlFull string,
	attrs ...ChangeCountAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			append(
				m.conv(attrs),
				attribute.String("vcs.change.state", string(changeState)),
				attribute.String("vcs.repository.url.full", repositoryUrlFull),
			)...,
		),
	)
}

func (m ChangeCount) conv(in []ChangeCountAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.changeCountAttr()
	}
	return out
}

// ChangeCountAttr is an optional attribute for the ChangeCount instrument.
type ChangeCountAttr interface {
    changeCountAttr() attribute.KeyValue
}

type changeCountAttr struct {
	kv attribute.KeyValue
}

func (a changeCountAttr) changeCountAttr() attribute.KeyValue {
    return a.kv
}

// RepositoryName returns an optional attribute for the "vcs.repository.name"
// semantic convention. It represents the human readable name of the repository.
// It SHOULD NOT include any additional identifier like Group/SubGroup in GitLab
// or organization in GitHub.
func (ChangeCount) RepositoryName(val string) ChangeCountAttr {
	return changeCountAttr{kv: attribute.String("vcs.repository.name", val)}
}

// VCSChangeDuration is an instrument used to record metric values conforming to
// the "vcs.change.duration" semantic conventions. It represents the time
// duration a change (pull request/merge request/changelist) has been in a given
// state.
type ChangeDuration struct {
	inst metric.Float64Gauge
}

// NewChangeDuration returns a new ChangeDuration instrument.
func NewChangeDuration(m metric.Meter) (ChangeDuration, error) {
	i, err := m.Float64Gauge(
	    "vcs.change.duration",
	    metric.WithDescription("The time duration a change (pull request/merge request/changelist) has been in a given state."),
	    metric.WithUnit("s"),
	)
	if err != nil {
	    return ChangeDuration{}, err
	}
	return ChangeDuration{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (ChangeDuration) Name() string {
	return "vcs.change.duration"
}

// Unit returns the semantic convention unit of the instrument
func (ChangeDuration) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (ChangeDuration) Description() string {
	return "The time duration a change (pull request/merge request/changelist) has been in a given state."
}

// Record records incr to the existing count.
//
// The vcsChangeState is the the state of the change (pull request/merge
// request/changelist).
//
// The vcsRefHeadName is the the name of the [reference] such as **branch** or
// **tag** in the repository.
//
// The vcsRepositoryUrlFull is the the [canonical URL] of the repository
// providing the complete HTTP(S) address in order to locate and identify the
// repository through a browser.
//
// All additional attrs passed are included in the recorded value.
//
// [reference]: https://git-scm.com/docs/gitglossary#def_ref
// [canonical URL]: https://support.google.com/webmasters/answer/10347851?hl=en#:~:text=A%20canonical%20URL%20is%20the,Google%20chooses%20one%20as%20canonical.
func (m ChangeDuration) Record(
    ctx context.Context,
    val float64,
	changeState ChangeStateAttr,
	refHeadName string,
	repositoryUrlFull string,
	attrs ...ChangeDurationAttr,
) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(
			append(
				m.conv(attrs),
				attribute.String("vcs.change.state", string(changeState)),
				attribute.String("vcs.ref.head.name", refHeadName),
				attribute.String("vcs.repository.url.full", repositoryUrlFull),
			)...,
		),
	)
}

func (m ChangeDuration) conv(in []ChangeDurationAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.changeDurationAttr()
	}
	return out
}

// ChangeDurationAttr is an optional attribute for the ChangeDuration instrument.
type ChangeDurationAttr interface {
    changeDurationAttr() attribute.KeyValue
}

type changeDurationAttr struct {
	kv attribute.KeyValue
}

func (a changeDurationAttr) changeDurationAttr() attribute.KeyValue {
    return a.kv
}

// RepositoryName returns an optional attribute for the "vcs.repository.name"
// semantic convention. It represents the human readable name of the repository.
// It SHOULD NOT include any additional identifier like Group/SubGroup in GitLab
// or organization in GitHub.
func (ChangeDuration) RepositoryName(val string) ChangeDurationAttr {
	return changeDurationAttr{kv: attribute.String("vcs.repository.name", val)}
}

// VCSChangeTimeToApproval is an instrument used to record metric values
// conforming to the "vcs.change.time_to_approval" semantic conventions. It
// represents the amount of time since its creation it took a change (pull
// request/merge request/changelist) to get the first approval.
type ChangeTimeToApproval struct {
	inst metric.Float64Gauge
}

// NewChangeTimeToApproval returns a new ChangeTimeToApproval instrument.
func NewChangeTimeToApproval(m metric.Meter) (ChangeTimeToApproval, error) {
	i, err := m.Float64Gauge(
	    "vcs.change.time_to_approval",
	    metric.WithDescription("The amount of time since its creation it took a change (pull request/merge request/changelist) to get the first approval."),
	    metric.WithUnit("s"),
	)
	if err != nil {
	    return ChangeTimeToApproval{}, err
	}
	return ChangeTimeToApproval{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (ChangeTimeToApproval) Name() string {
	return "vcs.change.time_to_approval"
}

// Unit returns the semantic convention unit of the instrument
func (ChangeTimeToApproval) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (ChangeTimeToApproval) Description() string {
	return "The amount of time since its creation it took a change (pull request/merge request/changelist) to get the first approval."
}

// Record records incr to the existing count.
//
// The vcsRefHeadName is the the name of the [reference] such as **branch** or
// **tag** in the repository.
//
// The vcsRepositoryUrlFull is the the [canonical URL] of the repository
// providing the complete HTTP(S) address in order to locate and identify the
// repository through a browser.
//
// All additional attrs passed are included in the recorded value.
//
// [reference]: https://git-scm.com/docs/gitglossary#def_ref
// [canonical URL]: https://support.google.com/webmasters/answer/10347851?hl=en#:~:text=A%20canonical%20URL%20is%20the,Google%20chooses%20one%20as%20canonical.
func (m ChangeTimeToApproval) Record(
    ctx context.Context,
    val float64,
	refHeadName string,
	repositoryUrlFull string,
	attrs ...ChangeTimeToApprovalAttr,
) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(
			append(
				m.conv(attrs),
				attribute.String("vcs.ref.head.name", refHeadName),
				attribute.String("vcs.repository.url.full", repositoryUrlFull),
			)...,
		),
	)
}

func (m ChangeTimeToApproval) conv(in []ChangeTimeToApprovalAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.changeTimeToApprovalAttr()
	}
	return out
}

// ChangeTimeToApprovalAttr is an optional attribute for the ChangeTimeToApproval
// instrument.
type ChangeTimeToApprovalAttr interface {
    changeTimeToApprovalAttr() attribute.KeyValue
}

type changeTimeToApprovalAttr struct {
	kv attribute.KeyValue
}

func (a changeTimeToApprovalAttr) changeTimeToApprovalAttr() attribute.KeyValue {
    return a.kv
}

// RefBaseName returns an optional attribute for the "vcs.ref.base.name" semantic
// convention. It represents the name of the [reference] such as **branch** or
// **tag** in the repository.
//
// [reference]: https://git-scm.com/docs/gitglossary#def_ref
func (ChangeTimeToApproval) RefBaseName(val string) ChangeTimeToApprovalAttr {
	return changeTimeToApprovalAttr{kv: attribute.String("vcs.ref.base.name", val)}
}

// RepositoryName returns an optional attribute for the "vcs.repository.name"
// semantic convention. It represents the human readable name of the repository.
// It SHOULD NOT include any additional identifier like Group/SubGroup in GitLab
// or organization in GitHub.
func (ChangeTimeToApproval) RepositoryName(val string) ChangeTimeToApprovalAttr {
	return changeTimeToApprovalAttr{kv: attribute.String("vcs.repository.name", val)}
}

// RefBaseRevision returns an optional attribute for the "vcs.ref.base.revision"
// semantic convention. It represents the revision, literally [revised version],
// The revision most often refers to a commit object in Git, or a revision number
// in SVN.
//
// [revised version]: https://www.merriam-webster.com/dictionary/revision
func (ChangeTimeToApproval) RefBaseRevision(val string) ChangeTimeToApprovalAttr {
	return changeTimeToApprovalAttr{kv: attribute.String("vcs.ref.base.revision", val)}
}

// RefHeadRevision returns an optional attribute for the "vcs.ref.head.revision"
// semantic convention. It represents the revision, literally [revised version],
// The revision most often refers to a commit object in Git, or a revision number
// in SVN.
//
// [revised version]: https://www.merriam-webster.com/dictionary/revision
func (ChangeTimeToApproval) RefHeadRevision(val string) ChangeTimeToApprovalAttr {
	return changeTimeToApprovalAttr{kv: attribute.String("vcs.ref.head.revision", val)}
}

// VCSChangeTimeToMerge is an instrument used to record metric values conforming
// to the "vcs.change.time_to_merge" semantic conventions. It represents the
// amount of time since its creation it took a change (pull request/merge
// request/changelist) to get merged into the target(base) ref.
type ChangeTimeToMerge struct {
	inst metric.Float64Gauge
}

// NewChangeTimeToMerge returns a new ChangeTimeToMerge instrument.
func NewChangeTimeToMerge(m metric.Meter) (ChangeTimeToMerge, error) {
	i, err := m.Float64Gauge(
	    "vcs.change.time_to_merge",
	    metric.WithDescription("The amount of time since its creation it took a change (pull request/merge request/changelist) to get merged into the target(base) ref."),
	    metric.WithUnit("s"),
	)
	if err != nil {
	    return ChangeTimeToMerge{}, err
	}
	return ChangeTimeToMerge{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (ChangeTimeToMerge) Name() string {
	return "vcs.change.time_to_merge"
}

// Unit returns the semantic convention unit of the instrument
func (ChangeTimeToMerge) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (ChangeTimeToMerge) Description() string {
	return "The amount of time since its creation it took a change (pull request/merge request/changelist) to get merged into the target(base) ref."
}

// Record records incr to the existing count.
//
// The vcsRefHeadName is the the name of the [reference] such as **branch** or
// **tag** in the repository.
//
// The vcsRepositoryUrlFull is the the [canonical URL] of the repository
// providing the complete HTTP(S) address in order to locate and identify the
// repository through a browser.
//
// All additional attrs passed are included in the recorded value.
//
// [reference]: https://git-scm.com/docs/gitglossary#def_ref
// [canonical URL]: https://support.google.com/webmasters/answer/10347851?hl=en#:~:text=A%20canonical%20URL%20is%20the,Google%20chooses%20one%20as%20canonical.
func (m ChangeTimeToMerge) Record(
    ctx context.Context,
    val float64,
	refHeadName string,
	repositoryUrlFull string,
	attrs ...ChangeTimeToMergeAttr,
) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(
			append(
				m.conv(attrs),
				attribute.String("vcs.ref.head.name", refHeadName),
				attribute.String("vcs.repository.url.full", repositoryUrlFull),
			)...,
		),
	)
}

func (m ChangeTimeToMerge) conv(in []ChangeTimeToMergeAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.changeTimeToMergeAttr()
	}
	return out
}

// ChangeTimeToMergeAttr is an optional attribute for the ChangeTimeToMerge
// instrument.
type ChangeTimeToMergeAttr interface {
    changeTimeToMergeAttr() attribute.KeyValue
}

type changeTimeToMergeAttr struct {
	kv attribute.KeyValue
}

func (a changeTimeToMergeAttr) changeTimeToMergeAttr() attribute.KeyValue {
    return a.kv
}

// RefBaseName returns an optional attribute for the "vcs.ref.base.name" semantic
// convention. It represents the name of the [reference] such as **branch** or
// **tag** in the repository.
//
// [reference]: https://git-scm.com/docs/gitglossary#def_ref
func (ChangeTimeToMerge) RefBaseName(val string) ChangeTimeToMergeAttr {
	return changeTimeToMergeAttr{kv: attribute.String("vcs.ref.base.name", val)}
}

// RepositoryName returns an optional attribute for the "vcs.repository.name"
// semantic convention. It represents the human readable name of the repository.
// It SHOULD NOT include any additional identifier like Group/SubGroup in GitLab
// or organization in GitHub.
func (ChangeTimeToMerge) RepositoryName(val string) ChangeTimeToMergeAttr {
	return changeTimeToMergeAttr{kv: attribute.String("vcs.repository.name", val)}
}

// RefBaseRevision returns an optional attribute for the "vcs.ref.base.revision"
// semantic convention. It represents the revision, literally [revised version],
// The revision most often refers to a commit object in Git, or a revision number
// in SVN.
//
// [revised version]: https://www.merriam-webster.com/dictionary/revision
func (ChangeTimeToMerge) RefBaseRevision(val string) ChangeTimeToMergeAttr {
	return changeTimeToMergeAttr{kv: attribute.String("vcs.ref.base.revision", val)}
}

// RefHeadRevision returns an optional attribute for the "vcs.ref.head.revision"
// semantic convention. It represents the revision, literally [revised version],
// The revision most often refers to a commit object in Git, or a revision number
// in SVN.
//
// [revised version]: https://www.merriam-webster.com/dictionary/revision
func (ChangeTimeToMerge) RefHeadRevision(val string) ChangeTimeToMergeAttr {
	return changeTimeToMergeAttr{kv: attribute.String("vcs.ref.head.revision", val)}
}

// VCSContributorCount is an instrument used to record metric values conforming
// to the "vcs.contributor.count" semantic conventions. It represents the number
// of unique contributors to a repository.
type ContributorCount struct {
	inst metric.Int64Gauge
}

// NewContributorCount returns a new ContributorCount instrument.
func NewContributorCount(m metric.Meter) (ContributorCount, error) {
	i, err := m.Int64Gauge(
	    "vcs.contributor.count",
	    metric.WithDescription("The number of unique contributors to a repository"),
	    metric.WithUnit("{contributor}"),
	)
	if err != nil {
	    return ContributorCount{}, err
	}
	return ContributorCount{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (ContributorCount) Name() string {
	return "vcs.contributor.count"
}

// Unit returns the semantic convention unit of the instrument
func (ContributorCount) Unit() string {
	return "{contributor}"
}

// Description returns the semantic convention description of the instrument
func (ContributorCount) Description() string {
	return "The number of unique contributors to a repository"
}

// Record records incr to the existing count.
//
// The vcsRepositoryUrlFull is the the [canonical URL] of the repository
// providing the complete HTTP(S) address in order to locate and identify the
// repository through a browser.
//
// All additional attrs passed are included in the recorded value.
//
// [canonical URL]: https://support.google.com/webmasters/answer/10347851?hl=en#:~:text=A%20canonical%20URL%20is%20the,Google%20chooses%20one%20as%20canonical.
func (m ContributorCount) Record(
    ctx context.Context,
    val int64,
	repositoryUrlFull string,
	attrs ...ContributorCountAttr,
) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(
			append(
				m.conv(attrs),
				attribute.String("vcs.repository.url.full", repositoryUrlFull),
			)...,
		),
	)
}

func (m ContributorCount) conv(in []ContributorCountAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.contributorCountAttr()
	}
	return out
}

// ContributorCountAttr is an optional attribute for the ContributorCount
// instrument.
type ContributorCountAttr interface {
    contributorCountAttr() attribute.KeyValue
}

type contributorCountAttr struct {
	kv attribute.KeyValue
}

func (a contributorCountAttr) contributorCountAttr() attribute.KeyValue {
    return a.kv
}

// RepositoryName returns an optional attribute for the "vcs.repository.name"
// semantic convention. It represents the human readable name of the repository.
// It SHOULD NOT include any additional identifier like Group/SubGroup in GitLab
// or organization in GitHub.
func (ContributorCount) RepositoryName(val string) ContributorCountAttr {
	return contributorCountAttr{kv: attribute.String("vcs.repository.name", val)}
}

// VCSRefCount is an instrument used to record metric values conforming to the
// "vcs.ref.count" semantic conventions. It represents the number of refs of type
// branch or tag in a repository.
type RefCount struct {
	inst metric.Int64UpDownCounter
}

// NewRefCount returns a new RefCount instrument.
func NewRefCount(m metric.Meter) (RefCount, error) {
	i, err := m.Int64UpDownCounter(
	    "vcs.ref.count",
	    metric.WithDescription("The number of refs of type branch or tag in a repository."),
	    metric.WithUnit("{ref}"),
	)
	if err != nil {
	    return RefCount{}, err
	}
	return RefCount{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (RefCount) Name() string {
	return "vcs.ref.count"
}

// Unit returns the semantic convention unit of the instrument
func (RefCount) Unit() string {
	return "{ref}"
}

// Description returns the semantic convention description of the instrument
func (RefCount) Description() string {
	return "The number of refs of type branch or tag in a repository."
}

// Add adds incr to the existing count.
//
// The vcsRefType is the the type of the [reference] in the repository.
//
// The vcsRepositoryUrlFull is the the [canonical URL] of the repository
// providing the complete HTTP(S) address in order to locate and identify the
// repository through a browser.
//
// All additional attrs passed are included in the recorded value.
//
// [reference]: https://git-scm.com/docs/gitglossary#def_ref
// [canonical URL]: https://support.google.com/webmasters/answer/10347851?hl=en#:~:text=A%20canonical%20URL%20is%20the,Google%20chooses%20one%20as%20canonical.
func (m RefCount) Add(
    ctx context.Context,
    incr int64,
	refType RefTypeAttr,
	repositoryUrlFull string,
	attrs ...RefCountAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			append(
				m.conv(attrs),
				attribute.String("vcs.ref.type", string(refType)),
				attribute.String("vcs.repository.url.full", repositoryUrlFull),
			)...,
		),
	)
}

func (m RefCount) conv(in []RefCountAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.refCountAttr()
	}
	return out
}

// RefCountAttr is an optional attribute for the RefCount instrument.
type RefCountAttr interface {
    refCountAttr() attribute.KeyValue
}

type refCountAttr struct {
	kv attribute.KeyValue
}

func (a refCountAttr) refCountAttr() attribute.KeyValue {
    return a.kv
}

// RepositoryName returns an optional attribute for the "vcs.repository.name"
// semantic convention. It represents the human readable name of the repository.
// It SHOULD NOT include any additional identifier like Group/SubGroup in GitLab
// or organization in GitHub.
func (RefCount) RepositoryName(val string) RefCountAttr {
	return refCountAttr{kv: attribute.String("vcs.repository.name", val)}
}

// VCSRefLinesDelta is an instrument used to record metric values conforming to
// the "vcs.ref.lines_delta" semantic conventions. It represents the number of
// lines added/removed in a ref (branch) relative to the ref from the
// `vcs.ref.base.name` attribute.
type RefLinesDelta struct {
	inst metric.Int64Gauge
}

// NewRefLinesDelta returns a new RefLinesDelta instrument.
func NewRefLinesDelta(m metric.Meter) (RefLinesDelta, error) {
	i, err := m.Int64Gauge(
	    "vcs.ref.lines_delta",
	    metric.WithDescription("The number of lines added/removed in a ref (branch) relative to the ref from the `vcs.ref.base.name` attribute."),
	    metric.WithUnit("{line}"),
	)
	if err != nil {
	    return RefLinesDelta{}, err
	}
	return RefLinesDelta{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (RefLinesDelta) Name() string {
	return "vcs.ref.lines_delta"
}

// Unit returns the semantic convention unit of the instrument
func (RefLinesDelta) Unit() string {
	return "{line}"
}

// Description returns the semantic convention description of the instrument
func (RefLinesDelta) Description() string {
	return "The number of lines added/removed in a ref (branch) relative to the ref from the `vcs.ref.base.name` attribute."
}

// Record records incr to the existing count.
//
// The vcsLineChangeType is the the type of line change being measured on a
// branch or change.
//
// The vcsRefBaseName is the the name of the [reference] such as **branch** or
// **tag** in the repository.
//
// The vcsRefBaseType is the the type of the [reference] in the repository.
//
// The vcsRefHeadName is the the name of the [reference] such as **branch** or
// **tag** in the repository.
//
// The vcsRefHeadType is the the type of the [reference] in the repository.
//
// The vcsRepositoryUrlFull is the the [canonical URL] of the repository
// providing the complete HTTP(S) address in order to locate and identify the
// repository through a browser.
//
// All additional attrs passed are included in the recorded value.
//
// [reference]: https://git-scm.com/docs/gitglossary#def_ref
// [reference]: https://git-scm.com/docs/gitglossary#def_ref
// [reference]: https://git-scm.com/docs/gitglossary#def_ref
// [reference]: https://git-scm.com/docs/gitglossary#def_ref
// [canonical URL]: https://support.google.com/webmasters/answer/10347851?hl=en#:~:text=A%20canonical%20URL%20is%20the,Google%20chooses%20one%20as%20canonical.
func (m RefLinesDelta) Record(
    ctx context.Context,
    val int64,
	lineChangeType LineChangeTypeAttr,
	refBaseName string,
	refBaseType RefBaseTypeAttr,
	refHeadName string,
	refHeadType RefHeadTypeAttr,
	repositoryUrlFull string,
	attrs ...RefLinesDeltaAttr,
) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(
			append(
				m.conv(attrs),
				attribute.String("vcs.line_change.type", string(lineChangeType)),
				attribute.String("vcs.ref.base.name", refBaseName),
				attribute.String("vcs.ref.base.type", string(refBaseType)),
				attribute.String("vcs.ref.head.name", refHeadName),
				attribute.String("vcs.ref.head.type", string(refHeadType)),
				attribute.String("vcs.repository.url.full", repositoryUrlFull),
			)...,
		),
	)
}

func (m RefLinesDelta) conv(in []RefLinesDeltaAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.refLinesDeltaAttr()
	}
	return out
}

// RefLinesDeltaAttr is an optional attribute for the RefLinesDelta instrument.
type RefLinesDeltaAttr interface {
    refLinesDeltaAttr() attribute.KeyValue
}

type refLinesDeltaAttr struct {
	kv attribute.KeyValue
}

func (a refLinesDeltaAttr) refLinesDeltaAttr() attribute.KeyValue {
    return a.kv
}

// ChangeID returns an optional attribute for the "vcs.change.id" semantic
// convention. It represents the ID of the change (pull request/merge
// request/changelist) if applicable. This is usually a unique (within
// repository) identifier generated by the VCS system.
func (RefLinesDelta) ChangeID(val string) RefLinesDeltaAttr {
	return refLinesDeltaAttr{kv: attribute.String("vcs.change.id", val)}
}

// RepositoryName returns an optional attribute for the "vcs.repository.name"
// semantic convention. It represents the human readable name of the repository.
// It SHOULD NOT include any additional identifier like Group/SubGroup in GitLab
// or organization in GitHub.
func (RefLinesDelta) RepositoryName(val string) RefLinesDeltaAttr {
	return refLinesDeltaAttr{kv: attribute.String("vcs.repository.name", val)}
}

// VCSRefRevisionsDelta is an instrument used to record metric values conforming
// to the "vcs.ref.revisions_delta" semantic conventions. It represents the
// number of revisions (commits) a ref (branch) is ahead/behind the branch from
// the `vcs.ref.base.name` attribute.
type RefRevisionsDelta struct {
	inst metric.Int64Gauge
}

// NewRefRevisionsDelta returns a new RefRevisionsDelta instrument.
func NewRefRevisionsDelta(m metric.Meter) (RefRevisionsDelta, error) {
	i, err := m.Int64Gauge(
	    "vcs.ref.revisions_delta",
	    metric.WithDescription("The number of revisions (commits) a ref (branch) is ahead/behind the branch from the `vcs.ref.base.name` attribute"),
	    metric.WithUnit("{revision}"),
	)
	if err != nil {
	    return RefRevisionsDelta{}, err
	}
	return RefRevisionsDelta{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (RefRevisionsDelta) Name() string {
	return "vcs.ref.revisions_delta"
}

// Unit returns the semantic convention unit of the instrument
func (RefRevisionsDelta) Unit() string {
	return "{revision}"
}

// Description returns the semantic convention description of the instrument
func (RefRevisionsDelta) Description() string {
	return "The number of revisions (commits) a ref (branch) is ahead/behind the branch from the `vcs.ref.base.name` attribute"
}

// Record records incr to the existing count.
//
// The vcsRefBaseName is the the name of the [reference] such as **branch** or
// **tag** in the repository.
//
// The vcsRefBaseType is the the type of the [reference] in the repository.
//
// The vcsRefHeadName is the the name of the [reference] such as **branch** or
// **tag** in the repository.
//
// The vcsRefHeadType is the the type of the [reference] in the repository.
//
// The vcsRepositoryUrlFull is the the [canonical URL] of the repository
// providing the complete HTTP(S) address in order to locate and identify the
// repository through a browser.
//
// The vcsRevisionDeltaDirection is the the type of revision comparison.
//
// All additional attrs passed are included in the recorded value.
//
// [reference]: https://git-scm.com/docs/gitglossary#def_ref
// [reference]: https://git-scm.com/docs/gitglossary#def_ref
// [reference]: https://git-scm.com/docs/gitglossary#def_ref
// [reference]: https://git-scm.com/docs/gitglossary#def_ref
// [canonical URL]: https://support.google.com/webmasters/answer/10347851?hl=en#:~:text=A%20canonical%20URL%20is%20the,Google%20chooses%20one%20as%20canonical.
func (m RefRevisionsDelta) Record(
    ctx context.Context,
    val int64,
	refBaseName string,
	refBaseType RefBaseTypeAttr,
	refHeadName string,
	refHeadType RefHeadTypeAttr,
	repositoryUrlFull string,
	revisionDeltaDirection RevisionDeltaDirectionAttr,
	attrs ...RefRevisionsDeltaAttr,
) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(
			append(
				m.conv(attrs),
				attribute.String("vcs.ref.base.name", refBaseName),
				attribute.String("vcs.ref.base.type", string(refBaseType)),
				attribute.String("vcs.ref.head.name", refHeadName),
				attribute.String("vcs.ref.head.type", string(refHeadType)),
				attribute.String("vcs.repository.url.full", repositoryUrlFull),
				attribute.String("vcs.revision_delta.direction", string(revisionDeltaDirection)),
			)...,
		),
	)
}

func (m RefRevisionsDelta) conv(in []RefRevisionsDeltaAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.refRevisionsDeltaAttr()
	}
	return out
}

// RefRevisionsDeltaAttr is an optional attribute for the RefRevisionsDelta
// instrument.
type RefRevisionsDeltaAttr interface {
    refRevisionsDeltaAttr() attribute.KeyValue
}

type refRevisionsDeltaAttr struct {
	kv attribute.KeyValue
}

func (a refRevisionsDeltaAttr) refRevisionsDeltaAttr() attribute.KeyValue {
    return a.kv
}

// ChangeID returns an optional attribute for the "vcs.change.id" semantic
// convention. It represents the ID of the change (pull request/merge
// request/changelist) if applicable. This is usually a unique (within
// repository) identifier generated by the VCS system.
func (RefRevisionsDelta) ChangeID(val string) RefRevisionsDeltaAttr {
	return refRevisionsDeltaAttr{kv: attribute.String("vcs.change.id", val)}
}

// RepositoryName returns an optional attribute for the "vcs.repository.name"
// semantic convention. It represents the human readable name of the repository.
// It SHOULD NOT include any additional identifier like Group/SubGroup in GitLab
// or organization in GitHub.
func (RefRevisionsDelta) RepositoryName(val string) RefRevisionsDeltaAttr {
	return refRevisionsDeltaAttr{kv: attribute.String("vcs.repository.name", val)}
}

// VCSRefTime is an instrument used to record metric values conforming to the
// "vcs.ref.time" semantic conventions. It represents the time a ref (branch)
// created from the default branch (trunk) has existed. The `ref.type` attribute
// will always be `branch`.
type RefTime struct {
	inst metric.Float64Gauge
}

// NewRefTime returns a new RefTime instrument.
func NewRefTime(m metric.Meter) (RefTime, error) {
	i, err := m.Float64Gauge(
	    "vcs.ref.time",
	    metric.WithDescription("Time a ref (branch) created from the default branch (trunk) has existed. The `ref.type` attribute will always be `branch`"),
	    metric.WithUnit("s"),
	)
	if err != nil {
	    return RefTime{}, err
	}
	return RefTime{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (RefTime) Name() string {
	return "vcs.ref.time"
}

// Unit returns the semantic convention unit of the instrument
func (RefTime) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (RefTime) Description() string {
	return "Time a ref (branch) created from the default branch (trunk) has existed. The `ref.type` attribute will always be `branch`"
}

// Record records incr to the existing count.
//
// The vcsRefHeadName is the the name of the [reference] such as **branch** or
// **tag** in the repository.
//
// The vcsRefHeadType is the the type of the [reference] in the repository.
//
// The vcsRepositoryUrlFull is the the [canonical URL] of the repository
// providing the complete HTTP(S) address in order to locate and identify the
// repository through a browser.
//
// All additional attrs passed are included in the recorded value.
//
// [reference]: https://git-scm.com/docs/gitglossary#def_ref
// [reference]: https://git-scm.com/docs/gitglossary#def_ref
// [canonical URL]: https://support.google.com/webmasters/answer/10347851?hl=en#:~:text=A%20canonical%20URL%20is%20the,Google%20chooses%20one%20as%20canonical.
func (m RefTime) Record(
    ctx context.Context,
    val float64,
	refHeadName string,
	refHeadType RefHeadTypeAttr,
	repositoryUrlFull string,
	attrs ...RefTimeAttr,
) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(
			append(
				m.conv(attrs),
				attribute.String("vcs.ref.head.name", refHeadName),
				attribute.String("vcs.ref.head.type", string(refHeadType)),
				attribute.String("vcs.repository.url.full", repositoryUrlFull),
			)...,
		),
	)
}

func (m RefTime) conv(in []RefTimeAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.refTimeAttr()
	}
	return out
}

// RefTimeAttr is an optional attribute for the RefTime instrument.
type RefTimeAttr interface {
    refTimeAttr() attribute.KeyValue
}

type refTimeAttr struct {
	kv attribute.KeyValue
}

func (a refTimeAttr) refTimeAttr() attribute.KeyValue {
    return a.kv
}

// RepositoryName returns an optional attribute for the "vcs.repository.name"
// semantic convention. It represents the human readable name of the repository.
// It SHOULD NOT include any additional identifier like Group/SubGroup in GitLab
// or organization in GitHub.
func (RefTime) RepositoryName(val string) RefTimeAttr {
	return refTimeAttr{kv: attribute.String("vcs.repository.name", val)}
}

// VCSRepositoryCount is an instrument used to record metric values conforming to
// the "vcs.repository.count" semantic conventions. It represents the number of
// repositories in an organization.
type RepositoryCount struct {
	inst metric.Int64UpDownCounter
}

// NewRepositoryCount returns a new RepositoryCount instrument.
func NewRepositoryCount(m metric.Meter) (RepositoryCount, error) {
	i, err := m.Int64UpDownCounter(
	    "vcs.repository.count",
	    metric.WithDescription("The number of repositories in an organization."),
	    metric.WithUnit("{repository}"),
	)
	if err != nil {
	    return RepositoryCount{}, err
	}
	return RepositoryCount{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (RepositoryCount) Name() string {
	return "vcs.repository.count"
}

// Unit returns the semantic convention unit of the instrument
func (RepositoryCount) Unit() string {
	return "{repository}"
}

// Description returns the semantic convention description of the instrument
func (RepositoryCount) Description() string {
	return "The number of repositories in an organization."
}

func (m RepositoryCount) Add(ctx context.Context, incr int64) {
    m.inst.Add(ctx, incr)
}