# [1.0.0](https://github.com/ogarciacar/kubecheck/compare/v0.0.2...v1.0.0) (2026-02-03)


### Code Refactoring

* rebrand module from kluster1 to kubecheck ([#42](https://github.com/ogarciacar/kubecheck/issues/42)) ([cf37407](https://github.com/ogarciacar/kubecheck/commit/cf37407daf3d9ac6107399c1d0426bcf91981712))


### Features

* **ci:** automate version releases with semantic-release ([#41](https://github.com/ogarciacar/kubecheck/issues/41)) ([4b1a347](https://github.com/ogarciacar/kubecheck/commit/4b1a347406c9af30927e7796a008e81adfffb881))


### BREAKING CHANGES

* Package name changed from kluster1 to kubecheck. Users must update their imports.

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>

* fix: update all kluster1 references to kubecheck

Replace remaining kluster1 references in:
- Code: kluster1. → kubecheck.
- Comments: Update package references
- Strings: Temp directory prefix

All tests now pass successfully.

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>
