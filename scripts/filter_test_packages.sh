#!/bin/bash

# List of packages to exclude (adjust as needed)
exclude_packages=(
  "github.com/beka-birhanu/task_manager_final/infrastructure/repo/task" 
  "github.com/beka-birhanu/task_manager_final/infrastructure/repo/user"
  "github.com/beka-birhanu/task_manager_final/api/errors"
  "github.com/beka-birhanu/task_manager_final/api/router"
  "github.com/beka-birhanu/task_manager_final/api/controllers/base"
  "github.com/beka-birhanu/task_manager_final/app/task/command/add/command"
  "github.com/beka-birhanu/task_manager_final/app/task/command/add/handler"
  "github.com/beka-birhanu/task_manager_final/app/task/command/delete/command"
  "github.com/beka-birhanu/task_manager_final/app/task/command/delete/handler"
  "github.com/beka-birhanu/task_manager_final/app/task/command/update/command"
  "github.com/beka-birhanu/task_manager_final/app/task/command/update/handler"
  "github.com/beka-birhanu/task_manager_final/app/task/query/get/command"
  "github.com/beka-birhanu/task_manager_final/app/task/query/get_all/command"
  "github.com/beka-birhanu/task_manager_final/app/user/auth/common"
  "github.com/beka-birhanu/task_manager_final/domain/errors"
  "github.com/beka-birhanu/task_manager_final/domain/i_hash"
  "github.com/beka-birhanu/task_manager_final/infrastructure/db"
  "github.com/beka-birhanu/task_manager_final/scripts"
  "github.com/beka-birhanu/task_manager_final/app/common/cqrs/command"
  "github.com/beka-birhanu/task_manager_final/app/common/cqrs/query/mocks"
  "github.com/beka-birhanu/task_manager_final/app/common/i_jwt"
  "github.com/beka-birhanu/task_manager_final/app/common/i_jwt/mock"
  "github.com/beka-birhanu/task_manager_final/app/common/i_repo"
  "github.com/beka-birhanu/task_manager_final/app/common/cqrs/command/mocks"
  "github.com/beka-birhanu/task_manager_final/app/common/cqrs/query"
  "github.com/beka-birhanu/task_manager_final/app/common/i_repo/mocks"
)

# Find all packages with .go files
all_packages=$(go list ./api/... ./app/... ./infrastructure/...)

# Filter out the excluded packages
for exclude in "${exclude_packages[@]}"; do
  all_packages=$(echo "$all_packages" | grep -v "^${exclude}$")
done

# Convert to a space-separated list
test_packages_list=$(echo "$all_packages" | tr '\n' ' ')

# Print the list of packages
echo $test_packages_list

