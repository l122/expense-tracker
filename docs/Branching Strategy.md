# Branching strategy

## Branch names

`main` - for deployment; no direct pushing of commits

`dev` - for merging all development, that is done in all other branches; no direct pushing of commits

`fix/<name-name>` - for fixing bugs, e.g. `fix/failing-tests`

`feat/<name-name>` - for new features (it can conatin bug fix), e.g. `feat/new-feature`

`refact/<name-name>` - for refactoring, e.g. `refact/unit-tests`

## Case scenarios

### A developer needs to add a new feature

- They pull the latest changes from the remote branch dev into the local repository

- They create a branch `feat/new-feature-name`

- Do all the development in this branch

- The development may contain a bug fix

- The development should not contain refactoring of the old code, unless it's part of the new feature

- After the work is done, publish the branch and create a pull requets back into the `dev` branch on GitHub

- Fix any comments or faling tests during PR check / review by pushing new commints

- Merge the PR and enjoy the success :)

### A developer needs to fix a bug only, no refactoring, no new feature

- They pull the latest changes from the remote branch dev into the local repository

- They create a branch `fix/bug-name`

- Do all the development in this branch

- After the work is done, publish the branch and create a pull requets back into the `dev` branch on GitHub

- Fix any comments or faling tests during PR check / review by pushing new commints

- Merge the PR and enjoy the success :)

### A developer needs to refactor existing code

- They pull the latest changes from the remote branch dev into the local repository

- They create a branch `refact/refactored-feature-name`

- Do all the development in this branch

- Do refactoring only, no bug fix, no new feature

- After the work is done, publish the branch and create a pull requets back into the `dev` branch on GitHub

- Fix any comments or faling tests during PR check / review by pushing new commints

- Merge the PR and enjoy the success :)