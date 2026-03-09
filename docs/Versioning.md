# Versioning

The version number uses the conventional format: 
`<major>.<minor>.<patch>` 

Example: `1.23.1`

The version tags are assing automatically using [Github Tag Bump](https://github.com/marketplace/actions/github-tag-bump)

## Rules to bump version

Include one of the following tags in the commit message that is part of a merge commit into `main`: 

 - `#major` Increments the major version (e.g., from v1.2.3 to v2.0.0).

 - `#minor`: Increments the minor version (e.g., from v1.2.3 to v1.3.0).

 - `#patch`: Increments the patch version (e.g., from v1.2.3 to v1.2.4).

 - `#none`: Skips the version bump entirely, overriding any DEFAULT_BUMP setting.

If none of the tags is specified, the default is `#minor`.

Example: 

`feat: update main #minor`