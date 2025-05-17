// Code generated from semantic convention specification. DO NOT EDIT.

// Package httpconv provides types and functionality for OpenTelemetry semantic
// conventions in the "vcs" namespace.
package vcsconv

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
)

var (
	addOptPool = &sync.Pool{New: func() any { return &[]metric.AddOption{} }}
	recOptPool = &sync.Pool{New: func() any { return &[]metric.RecordOption{} }}
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

// ProviderNameAttr is an attribute conforming to the vcs.provider.name semantic
// conventions. It represents the name of the version control system provider.
type ProviderNameAttr string

var (
	// ProviderNameGithub is the [GitHub].
	//
	// [GitHub]: https://github.com
	ProviderNameGithub ProviderNameAttr = "github"
	// ProviderNameGitlab is the [GitLab].
	//
	// [GitLab]: https://gitlab.com
	ProviderNameGitlab ProviderNameAttr = "gitlab"
	// ProviderNameGittea is the [Gitea].
	//
	// [Gitea]: https://gitea.io
	ProviderNameGittea ProviderNameAttr = "gittea"
	// ProviderNameBitbucket is the [Bitbucket].
	//
	// [Bitbucket]: https://bitbucket.org
	ProviderNameBitbucket ProviderNameAttr = "bitbucket"
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

// ChangeCount is an instrument used to record metric values conforming to the
// "vcs.change.count" semantic conventions. It represents the number of changes
// (pull requests/merge requests/changelists) in a repository, categorized by
// their state (e.g. open or merged).
type ChangeCount struct {
	metric.Int64UpDownCounter
}

// NewChangeCount returns a new ChangeCount instrument.
func NewChangeCount(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ChangeCount, error) {
	// Check if the meter is nil.
	if m == nil {
		return ChangeCount{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"vcs.change.count",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("The number of changes (pull requests/merge requests/changelists) in a repository, categorized by their state (e.g. open or merged)"),
			metric.WithUnit("{change}"),
		}, opt...)...,
	)
	if err != nil {
	    return ChangeCount{noop.Int64UpDownCounter{}}, err
	}
	return ChangeCount{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ChangeCount) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
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
// The changeState is the the state of the change (pull request/merge
// request/changelist).
//
// The repositoryUrlFull is the the [canonical URL] of the repository providing
// the complete HTTP(S) address in order to locate and identify the repository
// through a browser.
//
// All additional attrs passed are included in the recorded value.
//
// [canonical URL]: https://support.google.com/webmasters/answer/10347851?hl=en#:~:text=A%20canonical%20URL%20is%20the,Google%20chooses%20one%20as%20canonical.
func (m ChangeCount) Add(
	ctx context.Context,
	incr int64,
	changeState ChangeStateAttr,
	repositoryUrlFull string,
	attrs ...attribute.KeyValue,
) {
	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("vcs.change.state", string(changeState)),
				attribute.String("vcs.repository.url.full", repositoryUrlFull),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrOwnerName returns an optional attribute for the "vcs.owner.name" semantic
// convention. It represents the group owner within the version control system.
func (ChangeCount) AttrOwnerName(val string) attribute.KeyValue {
	return attribute.String("vcs.owner.name", val)
}

// AttrRepositoryName returns an optional attribute for the "vcs.repository.name"
// semantic convention. It represents the human readable name of the repository.
// It SHOULD NOT include any additional identifier like Group/SubGroup in GitLab
// or organization in GitHub.
func (ChangeCount) AttrRepositoryName(val string) attribute.KeyValue {
	return attribute.String("vcs.repository.name", val)
}

// AttrProviderName returns an optional attribute for the "vcs.provider.name"
// semantic convention. It represents the name of the version control system
// provider.
func (ChangeCount) AttrProviderName(val ProviderNameAttr) attribute.KeyValue {
	return attribute.String("vcs.provider.name", string(val))
}

// ChangeDuration is an instrument used to record metric values conforming to the
// "vcs.change.duration" semantic conventions. It represents the time duration a
// change (pull request/merge request/changelist) has been in a given state.
type ChangeDuration struct {
	metric.Float64Gauge
}

// NewChangeDuration returns a new ChangeDuration instrument.
func NewChangeDuration(
	m metric.Meter,
	opt ...metric.Float64GaugeOption,
) (ChangeDuration, error) {
	// Check if the meter is nil.
	if m == nil {
		return ChangeDuration{noop.Float64Gauge{}}, nil
	}

	i, err := m.Float64Gauge(
		"vcs.change.duration",
		append([]metric.Float64GaugeOption{
			metric.WithDescription("The time duration a change (pull request/merge request/changelist) has been in a given state."),
			metric.WithUnit("s"),
		}, opt...)...,
	)
	if err != nil {
	    return ChangeDuration{noop.Float64Gauge{}}, err
	}
	return ChangeDuration{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ChangeDuration) Inst() metric.Float64Gauge {
	return m.Float64Gauge
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

// Record records val to the current distribution.
//
// The changeState is the the state of the change (pull request/merge
// request/changelist).
//
// The refHeadName is the the name of the [reference] such as **branch** or
// **tag** in the repository.
//
// The repositoryUrlFull is the the [canonical URL] of the repository providing
// the complete HTTP(S) address in order to locate and identify the repository
// through a browser.
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
	attrs ...attribute.KeyValue,
) {
	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("vcs.change.state", string(changeState)),
				attribute.String("vcs.ref.head.name", refHeadName),
				attribute.String("vcs.repository.url.full", repositoryUrlFull),
			)...,
		),
	)

	m.Float64Gauge.Record(ctx, val, *o...)
}

// AttrOwnerName returns an optional attribute for the "vcs.owner.name" semantic
// convention. It represents the group owner within the version control system.
func (ChangeDuration) AttrOwnerName(val string) attribute.KeyValue {
	return attribute.String("vcs.owner.name", val)
}

// AttrRepositoryName returns an optional attribute for the "vcs.repository.name"
// semantic convention. It represents the human readable name of the repository.
// It SHOULD NOT include any additional identifier like Group/SubGroup in GitLab
// or organization in GitHub.
func (ChangeDuration) AttrRepositoryName(val string) attribute.KeyValue {
	return attribute.String("vcs.repository.name", val)
}

// AttrProviderName returns an optional attribute for the "vcs.provider.name"
// semantic convention. It represents the name of the version control system
// provider.
func (ChangeDuration) AttrProviderName(val ProviderNameAttr) attribute.KeyValue {
	return attribute.String("vcs.provider.name", string(val))
}

// ChangeTimeToApproval is an instrument used to record metric values conforming
// to the "vcs.change.time_to_approval" semantic conventions. It represents the
// amount of time since its creation it took a change (pull request/merge
// request/changelist) to get the first approval.
type ChangeTimeToApproval struct {
	metric.Float64Gauge
}

// NewChangeTimeToApproval returns a new ChangeTimeToApproval instrument.
func NewChangeTimeToApproval(
	m metric.Meter,
	opt ...metric.Float64GaugeOption,
) (ChangeTimeToApproval, error) {
	// Check if the meter is nil.
	if m == nil {
		return ChangeTimeToApproval{noop.Float64Gauge{}}, nil
	}

	i, err := m.Float64Gauge(
		"vcs.change.time_to_approval",
		append([]metric.Float64GaugeOption{
			metric.WithDescription("The amount of time since its creation it took a change (pull request/merge request/changelist) to get the first approval."),
			metric.WithUnit("s"),
		}, opt...)...,
	)
	if err != nil {
	    return ChangeTimeToApproval{noop.Float64Gauge{}}, err
	}
	return ChangeTimeToApproval{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ChangeTimeToApproval) Inst() metric.Float64Gauge {
	return m.Float64Gauge
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

// Record records val to the current distribution.
//
// The refHeadName is the the name of the [reference] such as **branch** or
// **tag** in the repository.
//
// The repositoryUrlFull is the the [canonical URL] of the repository providing
// the complete HTTP(S) address in order to locate and identify the repository
// through a browser.
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
	attrs ...attribute.KeyValue,
) {
	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("vcs.ref.head.name", refHeadName),
				attribute.String("vcs.repository.url.full", repositoryUrlFull),
			)...,
		),
	)

	m.Float64Gauge.Record(ctx, val, *o...)
}

// AttrOwnerName returns an optional attribute for the "vcs.owner.name" semantic
// convention. It represents the group owner within the version control system.
func (ChangeTimeToApproval) AttrOwnerName(val string) attribute.KeyValue {
	return attribute.String("vcs.owner.name", val)
}

// AttrRefBaseName returns an optional attribute for the "vcs.ref.base.name"
// semantic convention. It represents the name of the [reference] such as
// **branch** or **tag** in the repository.
//
// [reference]: https://git-scm.com/docs/gitglossary#def_ref
func (ChangeTimeToApproval) AttrRefBaseName(val string) attribute.KeyValue {
	return attribute.String("vcs.ref.base.name", val)
}

// AttrRepositoryName returns an optional attribute for the "vcs.repository.name"
// semantic convention. It represents the human readable name of the repository.
// It SHOULD NOT include any additional identifier like Group/SubGroup in GitLab
// or organization in GitHub.
func (ChangeTimeToApproval) AttrRepositoryName(val string) attribute.KeyValue {
	return attribute.String("vcs.repository.name", val)
}

// AttrProviderName returns an optional attribute for the "vcs.provider.name"
// semantic convention. It represents the name of the version control system
// provider.
func (ChangeTimeToApproval) AttrProviderName(val ProviderNameAttr) attribute.KeyValue {
	return attribute.String("vcs.provider.name", string(val))
}

// AttrRefBaseRevision returns an optional attribute for the
// "vcs.ref.base.revision" semantic convention. It represents the revision,
// literally [revised version], The revision most often refers to a commit object
// in Git, or a revision number in SVN.
//
// [revised version]: https://www.merriam-webster.com/dictionary/revision
func (ChangeTimeToApproval) AttrRefBaseRevision(val string) attribute.KeyValue {
	return attribute.String("vcs.ref.base.revision", val)
}

// AttrRefHeadRevision returns an optional attribute for the
// "vcs.ref.head.revision" semantic convention. It represents the revision,
// literally [revised version], The revision most often refers to a commit object
// in Git, or a revision number in SVN.
//
// [revised version]: https://www.merriam-webster.com/dictionary/revision
func (ChangeTimeToApproval) AttrRefHeadRevision(val string) attribute.KeyValue {
	return attribute.String("vcs.ref.head.revision", val)
}

// ChangeTimeToMerge is an instrument used to record metric values conforming to
// the "vcs.change.time_to_merge" semantic conventions. It represents the amount
// of time since its creation it took a change (pull request/merge
// request/changelist) to get merged into the target(base) ref.
type ChangeTimeToMerge struct {
	metric.Float64Gauge
}

// NewChangeTimeToMerge returns a new ChangeTimeToMerge instrument.
func NewChangeTimeToMerge(
	m metric.Meter,
	opt ...metric.Float64GaugeOption,
) (ChangeTimeToMerge, error) {
	// Check if the meter is nil.
	if m == nil {
		return ChangeTimeToMerge{noop.Float64Gauge{}}, nil
	}

	i, err := m.Float64Gauge(
		"vcs.change.time_to_merge",
		append([]metric.Float64GaugeOption{
			metric.WithDescription("The amount of time since its creation it took a change (pull request/merge request/changelist) to get merged into the target(base) ref."),
			metric.WithUnit("s"),
		}, opt...)...,
	)
	if err != nil {
	    return ChangeTimeToMerge{noop.Float64Gauge{}}, err
	}
	return ChangeTimeToMerge{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ChangeTimeToMerge) Inst() metric.Float64Gauge {
	return m.Float64Gauge
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

// Record records val to the current distribution.
//
// The refHeadName is the the name of the [reference] such as **branch** or
// **tag** in the repository.
//
// The repositoryUrlFull is the the [canonical URL] of the repository providing
// the complete HTTP(S) address in order to locate and identify the repository
// through a browser.
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
	attrs ...attribute.KeyValue,
) {
	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("vcs.ref.head.name", refHeadName),
				attribute.String("vcs.repository.url.full", repositoryUrlFull),
			)...,
		),
	)

	m.Float64Gauge.Record(ctx, val, *o...)
}

// AttrOwnerName returns an optional attribute for the "vcs.owner.name" semantic
// convention. It represents the group owner within the version control system.
func (ChangeTimeToMerge) AttrOwnerName(val string) attribute.KeyValue {
	return attribute.String("vcs.owner.name", val)
}

// AttrRefBaseName returns an optional attribute for the "vcs.ref.base.name"
// semantic convention. It represents the name of the [reference] such as
// **branch** or **tag** in the repository.
//
// [reference]: https://git-scm.com/docs/gitglossary#def_ref
func (ChangeTimeToMerge) AttrRefBaseName(val string) attribute.KeyValue {
	return attribute.String("vcs.ref.base.name", val)
}

// AttrRepositoryName returns an optional attribute for the "vcs.repository.name"
// semantic convention. It represents the human readable name of the repository.
// It SHOULD NOT include any additional identifier like Group/SubGroup in GitLab
// or organization in GitHub.
func (ChangeTimeToMerge) AttrRepositoryName(val string) attribute.KeyValue {
	return attribute.String("vcs.repository.name", val)
}

// AttrProviderName returns an optional attribute for the "vcs.provider.name"
// semantic convention. It represents the name of the version control system
// provider.
func (ChangeTimeToMerge) AttrProviderName(val ProviderNameAttr) attribute.KeyValue {
	return attribute.String("vcs.provider.name", string(val))
}

// AttrRefBaseRevision returns an optional attribute for the
// "vcs.ref.base.revision" semantic convention. It represents the revision,
// literally [revised version], The revision most often refers to a commit object
// in Git, or a revision number in SVN.
//
// [revised version]: https://www.merriam-webster.com/dictionary/revision
func (ChangeTimeToMerge) AttrRefBaseRevision(val string) attribute.KeyValue {
	return attribute.String("vcs.ref.base.revision", val)
}

// AttrRefHeadRevision returns an optional attribute for the
// "vcs.ref.head.revision" semantic convention. It represents the revision,
// literally [revised version], The revision most often refers to a commit object
// in Git, or a revision number in SVN.
//
// [revised version]: https://www.merriam-webster.com/dictionary/revision
func (ChangeTimeToMerge) AttrRefHeadRevision(val string) attribute.KeyValue {
	return attribute.String("vcs.ref.head.revision", val)
}

// ContributorCount is an instrument used to record metric values conforming to
// the "vcs.contributor.count" semantic conventions. It represents the number of
// unique contributors to a repository.
type ContributorCount struct {
	metric.Int64Gauge
}

// NewContributorCount returns a new ContributorCount instrument.
func NewContributorCount(
	m metric.Meter,
	opt ...metric.Int64GaugeOption,
) (ContributorCount, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContributorCount{noop.Int64Gauge{}}, nil
	}

	i, err := m.Int64Gauge(
		"vcs.contributor.count",
		append([]metric.Int64GaugeOption{
			metric.WithDescription("The number of unique contributors to a repository"),
			metric.WithUnit("{contributor}"),
		}, opt...)...,
	)
	if err != nil {
	    return ContributorCount{noop.Int64Gauge{}}, err
	}
	return ContributorCount{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ContributorCount) Inst() metric.Int64Gauge {
	return m.Int64Gauge
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

// Record records val to the current distribution.
//
// The repositoryUrlFull is the the [canonical URL] of the repository providing
// the complete HTTP(S) address in order to locate and identify the repository
// through a browser.
//
// All additional attrs passed are included in the recorded value.
//
// [canonical URL]: https://support.google.com/webmasters/answer/10347851?hl=en#:~:text=A%20canonical%20URL%20is%20the,Google%20chooses%20one%20as%20canonical.
func (m ContributorCount) Record(
	ctx context.Context,
	val int64,
	repositoryUrlFull string,
	attrs ...attribute.KeyValue,
) {
	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("vcs.repository.url.full", repositoryUrlFull),
			)...,
		),
	)

	m.Int64Gauge.Record(ctx, val, *o...)
}

// AttrOwnerName returns an optional attribute for the "vcs.owner.name" semantic
// convention. It represents the group owner within the version control system.
func (ContributorCount) AttrOwnerName(val string) attribute.KeyValue {
	return attribute.String("vcs.owner.name", val)
}

// AttrRepositoryName returns an optional attribute for the "vcs.repository.name"
// semantic convention. It represents the human readable name of the repository.
// It SHOULD NOT include any additional identifier like Group/SubGroup in GitLab
// or organization in GitHub.
func (ContributorCount) AttrRepositoryName(val string) attribute.KeyValue {
	return attribute.String("vcs.repository.name", val)
}

// AttrProviderName returns an optional attribute for the "vcs.provider.name"
// semantic convention. It represents the name of the version control system
// provider.
func (ContributorCount) AttrProviderName(val ProviderNameAttr) attribute.KeyValue {
	return attribute.String("vcs.provider.name", string(val))
}

// RefCount is an instrument used to record metric values conforming to the
// "vcs.ref.count" semantic conventions. It represents the number of refs of type
// branch or tag in a repository.
type RefCount struct {
	metric.Int64UpDownCounter
}

// NewRefCount returns a new RefCount instrument.
func NewRefCount(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (RefCount, error) {
	// Check if the meter is nil.
	if m == nil {
		return RefCount{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"vcs.ref.count",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("The number of refs of type branch or tag in a repository."),
			metric.WithUnit("{ref}"),
		}, opt...)...,
	)
	if err != nil {
	    return RefCount{noop.Int64UpDownCounter{}}, err
	}
	return RefCount{i}, nil
}

// Inst returns the underlying metric instrument.
func (m RefCount) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
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
// The refType is the the type of the [reference] in the repository.
//
// The repositoryUrlFull is the the [canonical URL] of the repository providing
// the complete HTTP(S) address in order to locate and identify the repository
// through a browser.
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
	attrs ...attribute.KeyValue,
) {
	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("vcs.ref.type", string(refType)),
				attribute.String("vcs.repository.url.full", repositoryUrlFull),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrOwnerName returns an optional attribute for the "vcs.owner.name" semantic
// convention. It represents the group owner within the version control system.
func (RefCount) AttrOwnerName(val string) attribute.KeyValue {
	return attribute.String("vcs.owner.name", val)
}

// AttrRepositoryName returns an optional attribute for the "vcs.repository.name"
// semantic convention. It represents the human readable name of the repository.
// It SHOULD NOT include any additional identifier like Group/SubGroup in GitLab
// or organization in GitHub.
func (RefCount) AttrRepositoryName(val string) attribute.KeyValue {
	return attribute.String("vcs.repository.name", val)
}

// AttrProviderName returns an optional attribute for the "vcs.provider.name"
// semantic convention. It represents the name of the version control system
// provider.
func (RefCount) AttrProviderName(val ProviderNameAttr) attribute.KeyValue {
	return attribute.String("vcs.provider.name", string(val))
}

// RefLinesDelta is an instrument used to record metric values conforming to the
// "vcs.ref.lines_delta" semantic conventions. It represents the number of lines
// added/removed in a ref (branch) relative to the ref from the
// `vcs.ref.base.name` attribute.
type RefLinesDelta struct {
	metric.Int64Gauge
}

// NewRefLinesDelta returns a new RefLinesDelta instrument.
func NewRefLinesDelta(
	m metric.Meter,
	opt ...metric.Int64GaugeOption,
) (RefLinesDelta, error) {
	// Check if the meter is nil.
	if m == nil {
		return RefLinesDelta{noop.Int64Gauge{}}, nil
	}

	i, err := m.Int64Gauge(
		"vcs.ref.lines_delta",
		append([]metric.Int64GaugeOption{
			metric.WithDescription("The number of lines added/removed in a ref (branch) relative to the ref from the `vcs.ref.base.name` attribute."),
			metric.WithUnit("{line}"),
		}, opt...)...,
	)
	if err != nil {
	    return RefLinesDelta{noop.Int64Gauge{}}, err
	}
	return RefLinesDelta{i}, nil
}

// Inst returns the underlying metric instrument.
func (m RefLinesDelta) Inst() metric.Int64Gauge {
	return m.Int64Gauge
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

// Record records val to the current distribution.
//
// The lineChangeType is the the type of line change being measured on a branch
// or change.
//
// The refBaseName is the the name of the [reference] such as **branch** or
// **tag** in the repository.
//
// The refBaseType is the the type of the [reference] in the repository.
//
// The refHeadName is the the name of the [reference] such as **branch** or
// **tag** in the repository.
//
// The refHeadType is the the type of the [reference] in the repository.
//
// The repositoryUrlFull is the the [canonical URL] of the repository providing
// the complete HTTP(S) address in order to locate and identify the repository
// through a browser.
//
// All additional attrs passed are included in the recorded value.
//
// [reference]: https://git-scm.com/docs/gitglossary#def_ref
// [reference]: https://git-scm.com/docs/gitglossary#def_ref
// [reference]: https://git-scm.com/docs/gitglossary#def_ref
// [reference]: https://git-scm.com/docs/gitglossary#def_ref
// [canonical URL]: https://support.google.com/webmasters/answer/10347851?hl=en#:~:text=A%20canonical%20URL%20is%20the,Google%20chooses%20one%20as%20canonical.
//
// This metric should be reported for each `vcs.line_change.type` value. For
// example if a ref added 3 lines and removed 2 lines,
// instrumentation SHOULD report two measurements: 3 and 2 (both positive
// numbers).
// If number of lines added/removed should be calculated from the start of time,
// then `vcs.ref.base.name` SHOULD be set to an empty string.
func (m RefLinesDelta) Record(
	ctx context.Context,
	val int64,
	lineChangeType LineChangeTypeAttr,
	refBaseName string,
	refBaseType RefBaseTypeAttr,
	refHeadName string,
	refHeadType RefHeadTypeAttr,
	repositoryUrlFull string,
	attrs ...attribute.KeyValue,
) {
	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("vcs.line_change.type", string(lineChangeType)),
				attribute.String("vcs.ref.base.name", refBaseName),
				attribute.String("vcs.ref.base.type", string(refBaseType)),
				attribute.String("vcs.ref.head.name", refHeadName),
				attribute.String("vcs.ref.head.type", string(refHeadType)),
				attribute.String("vcs.repository.url.full", repositoryUrlFull),
			)...,
		),
	)

	m.Int64Gauge.Record(ctx, val, *o...)
}

// AttrChangeID returns an optional attribute for the "vcs.change.id" semantic
// convention. It represents the ID of the change (pull request/merge
// request/changelist) if applicable. This is usually a unique (within
// repository) identifier generated by the VCS system.
func (RefLinesDelta) AttrChangeID(val string) attribute.KeyValue {
	return attribute.String("vcs.change.id", val)
}

// AttrOwnerName returns an optional attribute for the "vcs.owner.name" semantic
// convention. It represents the group owner within the version control system.
func (RefLinesDelta) AttrOwnerName(val string) attribute.KeyValue {
	return attribute.String("vcs.owner.name", val)
}

// AttrRepositoryName returns an optional attribute for the "vcs.repository.name"
// semantic convention. It represents the human readable name of the repository.
// It SHOULD NOT include any additional identifier like Group/SubGroup in GitLab
// or organization in GitHub.
func (RefLinesDelta) AttrRepositoryName(val string) attribute.KeyValue {
	return attribute.String("vcs.repository.name", val)
}

// AttrProviderName returns an optional attribute for the "vcs.provider.name"
// semantic convention. It represents the name of the version control system
// provider.
func (RefLinesDelta) AttrProviderName(val ProviderNameAttr) attribute.KeyValue {
	return attribute.String("vcs.provider.name", string(val))
}

// RefRevisionsDelta is an instrument used to record metric values conforming to
// the "vcs.ref.revisions_delta" semantic conventions. It represents the number
// of revisions (commits) a ref (branch) is ahead/behind the branch from the
// `vcs.ref.base.name` attribute.
type RefRevisionsDelta struct {
	metric.Int64Gauge
}

// NewRefRevisionsDelta returns a new RefRevisionsDelta instrument.
func NewRefRevisionsDelta(
	m metric.Meter,
	opt ...metric.Int64GaugeOption,
) (RefRevisionsDelta, error) {
	// Check if the meter is nil.
	if m == nil {
		return RefRevisionsDelta{noop.Int64Gauge{}}, nil
	}

	i, err := m.Int64Gauge(
		"vcs.ref.revisions_delta",
		append([]metric.Int64GaugeOption{
			metric.WithDescription("The number of revisions (commits) a ref (branch) is ahead/behind the branch from the `vcs.ref.base.name` attribute"),
			metric.WithUnit("{revision}"),
		}, opt...)...,
	)
	if err != nil {
	    return RefRevisionsDelta{noop.Int64Gauge{}}, err
	}
	return RefRevisionsDelta{i}, nil
}

// Inst returns the underlying metric instrument.
func (m RefRevisionsDelta) Inst() metric.Int64Gauge {
	return m.Int64Gauge
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

// Record records val to the current distribution.
//
// The refBaseName is the the name of the [reference] such as **branch** or
// **tag** in the repository.
//
// The refBaseType is the the type of the [reference] in the repository.
//
// The refHeadName is the the name of the [reference] such as **branch** or
// **tag** in the repository.
//
// The refHeadType is the the type of the [reference] in the repository.
//
// The repositoryUrlFull is the the [canonical URL] of the repository providing
// the complete HTTP(S) address in order to locate and identify the repository
// through a browser.
//
// The revisionDeltaDirection is the the type of revision comparison.
//
// All additional attrs passed are included in the recorded value.
//
// [reference]: https://git-scm.com/docs/gitglossary#def_ref
// [reference]: https://git-scm.com/docs/gitglossary#def_ref
// [reference]: https://git-scm.com/docs/gitglossary#def_ref
// [reference]: https://git-scm.com/docs/gitglossary#def_ref
// [canonical URL]: https://support.google.com/webmasters/answer/10347851?hl=en#:~:text=A%20canonical%20URL%20is%20the,Google%20chooses%20one%20as%20canonical.
//
// This metric should be reported for each `vcs.revision_delta.direction` value.
// For example if branch `a` is 3 commits behind and 2 commits ahead of `trunk`,
// instrumentation SHOULD report two measurements: 3 and 2 (both positive
// numbers) and `vcs.ref.base.name` is set to `trunk`.
func (m RefRevisionsDelta) Record(
	ctx context.Context,
	val int64,
	refBaseName string,
	refBaseType RefBaseTypeAttr,
	refHeadName string,
	refHeadType RefHeadTypeAttr,
	repositoryUrlFull string,
	revisionDeltaDirection RevisionDeltaDirectionAttr,
	attrs ...attribute.KeyValue,
) {
	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("vcs.ref.base.name", refBaseName),
				attribute.String("vcs.ref.base.type", string(refBaseType)),
				attribute.String("vcs.ref.head.name", refHeadName),
				attribute.String("vcs.ref.head.type", string(refHeadType)),
				attribute.String("vcs.repository.url.full", repositoryUrlFull),
				attribute.String("vcs.revision_delta.direction", string(revisionDeltaDirection)),
			)...,
		),
	)

	m.Int64Gauge.Record(ctx, val, *o...)
}

// AttrChangeID returns an optional attribute for the "vcs.change.id" semantic
// convention. It represents the ID of the change (pull request/merge
// request/changelist) if applicable. This is usually a unique (within
// repository) identifier generated by the VCS system.
func (RefRevisionsDelta) AttrChangeID(val string) attribute.KeyValue {
	return attribute.String("vcs.change.id", val)
}

// AttrOwnerName returns an optional attribute for the "vcs.owner.name" semantic
// convention. It represents the group owner within the version control system.
func (RefRevisionsDelta) AttrOwnerName(val string) attribute.KeyValue {
	return attribute.String("vcs.owner.name", val)
}

// AttrRepositoryName returns an optional attribute for the "vcs.repository.name"
// semantic convention. It represents the human readable name of the repository.
// It SHOULD NOT include any additional identifier like Group/SubGroup in GitLab
// or organization in GitHub.
func (RefRevisionsDelta) AttrRepositoryName(val string) attribute.KeyValue {
	return attribute.String("vcs.repository.name", val)
}

// AttrProviderName returns an optional attribute for the "vcs.provider.name"
// semantic convention. It represents the name of the version control system
// provider.
func (RefRevisionsDelta) AttrProviderName(val ProviderNameAttr) attribute.KeyValue {
	return attribute.String("vcs.provider.name", string(val))
}

// RefTime is an instrument used to record metric values conforming to the
// "vcs.ref.time" semantic conventions. It represents the time a ref (branch)
// created from the default branch (trunk) has existed. The `ref.type` attribute
// will always be `branch`.
type RefTime struct {
	metric.Float64Gauge
}

// NewRefTime returns a new RefTime instrument.
func NewRefTime(
	m metric.Meter,
	opt ...metric.Float64GaugeOption,
) (RefTime, error) {
	// Check if the meter is nil.
	if m == nil {
		return RefTime{noop.Float64Gauge{}}, nil
	}

	i, err := m.Float64Gauge(
		"vcs.ref.time",
		append([]metric.Float64GaugeOption{
			metric.WithDescription("Time a ref (branch) created from the default branch (trunk) has existed. The `ref.type` attribute will always be `branch`"),
			metric.WithUnit("s"),
		}, opt...)...,
	)
	if err != nil {
	    return RefTime{noop.Float64Gauge{}}, err
	}
	return RefTime{i}, nil
}

// Inst returns the underlying metric instrument.
func (m RefTime) Inst() metric.Float64Gauge {
	return m.Float64Gauge
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

// Record records val to the current distribution.
//
// The refHeadName is the the name of the [reference] such as **branch** or
// **tag** in the repository.
//
// The refHeadType is the the type of the [reference] in the repository.
//
// The repositoryUrlFull is the the [canonical URL] of the repository providing
// the complete HTTP(S) address in order to locate and identify the repository
// through a browser.
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
	attrs ...attribute.KeyValue,
) {
	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("vcs.ref.head.name", refHeadName),
				attribute.String("vcs.ref.head.type", string(refHeadType)),
				attribute.String("vcs.repository.url.full", repositoryUrlFull),
			)...,
		),
	)

	m.Float64Gauge.Record(ctx, val, *o...)
}

// AttrOwnerName returns an optional attribute for the "vcs.owner.name" semantic
// convention. It represents the group owner within the version control system.
func (RefTime) AttrOwnerName(val string) attribute.KeyValue {
	return attribute.String("vcs.owner.name", val)
}

// AttrRepositoryName returns an optional attribute for the "vcs.repository.name"
// semantic convention. It represents the human readable name of the repository.
// It SHOULD NOT include any additional identifier like Group/SubGroup in GitLab
// or organization in GitHub.
func (RefTime) AttrRepositoryName(val string) attribute.KeyValue {
	return attribute.String("vcs.repository.name", val)
}

// AttrProviderName returns an optional attribute for the "vcs.provider.name"
// semantic convention. It represents the name of the version control system
// provider.
func (RefTime) AttrProviderName(val ProviderNameAttr) attribute.KeyValue {
	return attribute.String("vcs.provider.name", string(val))
}

// RepositoryCount is an instrument used to record metric values conforming to
// the "vcs.repository.count" semantic conventions. It represents the number of
// repositories in an organization.
type RepositoryCount struct {
	metric.Int64UpDownCounter
}

// NewRepositoryCount returns a new RepositoryCount instrument.
func NewRepositoryCount(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (RepositoryCount, error) {
	// Check if the meter is nil.
	if m == nil {
		return RepositoryCount{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"vcs.repository.count",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("The number of repositories in an organization."),
			metric.WithUnit("{repository}"),
		}, opt...)...,
	)
	if err != nil {
	    return RepositoryCount{noop.Int64UpDownCounter{}}, err
	}
	return RepositoryCount{i}, nil
}

// Inst returns the underlying metric instrument.
func (m RepositoryCount) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
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

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
func (m RepositoryCount) Add(
	ctx context.Context,
	incr int64,
	attrs ...attribute.KeyValue,
) {
	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			attrs...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrOwnerName returns an optional attribute for the "vcs.owner.name" semantic
// convention. It represents the group owner within the version control system.
func (RepositoryCount) AttrOwnerName(val string) attribute.KeyValue {
	return attribute.String("vcs.owner.name", val)
}

// AttrProviderName returns an optional attribute for the "vcs.provider.name"
// semantic convention. It represents the name of the version control system
// provider.
func (RepositoryCount) AttrProviderName(val ProviderNameAttr) attribute.KeyValue {
	return attribute.String("vcs.provider.name", string(val))
}