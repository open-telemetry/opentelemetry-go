# Release Process

## Pre-Release
Update go.mod for submodules to depend on the new release which will happen
in the next step. This will create build failure for those who depend
on the master branch instead or released version. But they shouldn't be
depending on the master. So it is not a concern.

1. Run the pre-release script. It creates a branch pre_release_<tag> to make the changes.

    ```
    ./pre_release.sh -t <new tag>
    ```

2. Verify the changes.

    ```
    git diff master
    ```

3. Push the changes to upstream.

    ```
    git push
    ```

4. Create a PR on github and merge the PR once approved.


## Tag
Now create a new Tag on the commit hash of the changes made in pre-release step.
Use the same tag as used in the pre-release step.

1. Run the tag.sh script.

    ```
    ./tag.sh <new tag> <commit-hash>
    ```

2. Push tags upstream. Make sure you run this for all sub-modules as well.

    ```
    git push upstream <new tag>
    git push upstream <submodules-path/new tag>
    ...
    ```

## Release
Now create a release for the new `<new tag>` on github.
The release body should include all the release notes in the Changelog for this release.
Additionally, the `tag.sh` script generates commit logs since last release which can be used to suppliment the release notes.

## Verify Examples
After releasing run following script to verify that examples build outside of the otel repo.
The script copies examples into a different directory and builds them.

```
./verify_examples.sh
```

