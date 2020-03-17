# Release Process

## Pre-Release
Update go.mod for submodules to depend on the new release which will happen
in the next step. This will create build failure for those who depend
on the master branch instead or released version. But they shouldn't be
depending on the master. So it is not a concern.

- Run
```
./pre-release.sh -t <new tag>
```

- Verify the changes
```
git diff master
```

- Push changes to upstream
```
git push
```

- Create PR on github and merge it once approved.

## Tag
Now create a new Tag on the commit hash of the changes made in pre-release step.
Use the same tag as used in the pre-release step.

- Run
```
./tag.sh -t <new tag> -c <commit-hash>
```

- Push tags upstream. Make sure you run this for all sub-modules as well.
```
git push upstream <new tag>
git push upstream <submodules-path/new tag>
```

## Release
Now create a release for the new tag on github.
