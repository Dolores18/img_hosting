# 图片托管系统 API 文档

## 目录

1. [认证相关](#认证相关)
2. [用户管理](#用户管理)
3. [图片管理](#图片管理)
4. [标签管理](#标签管理)
5. [私有文件管理](#私有文件管理)
6. [令牌管理](#令牌管理)
7. [权限管理](#权限管理)

## 认证相关

### 用户注册

- **URL**: `/auth/register`
- **方法**: `POST`
- **请求体**:
  ```json
  {
    "name": "用户名",
    "email": "邮箱地址",// 必须是唯一的
    "psd": "密码", // 注意这里使用psd而不是password
    "age": 25
  }
  ```
- **响应**:
  ```json
  {
    "message": "注册成功",
    "user_id": 1
  }
  ```
- **错误响应**:
  ```json
  {
    "error": "邮箱已被使用"
  }
  ```

### 用户登录

- **URL**: `/auth/login`
- **方法**: `POST`
- **请求体**:
  ```json
  {
    "email": "邮箱地址",
    "password": "密码"
  }
  ```
- **响应**:
  ```json
  {
    "message": "登录成功",
    "token": "JWT令牌",
    "user_id": 1,
    "user_name": "用户名"
  }
  ```

## 用户管理

### 获取用户个人资料

- **URL**: `/users/profile`
- **方法**: `GET`
- **请求头**: `Authorization: Bearer {token}`
- **响应**:
  ```json
  {
    "user": {
      "id": 1,
      "name": "用户名",
      "email": "邮箱地址",
      "created_at": "创建时间",
      "status": "状态"
    }
  }
  ```

### 更新用户个人资料

- **URL**: `/users/profile`
- **方法**: `PUT`
- **请求头**: `Authorization: Bearer {token}`
- **请求体**:
  ```json
  {
    "name": "新用户名",
    "email": "新邮箱地址",
    "password": "新密码"
  }
  ```
- **响应**:
  ```json
  {
    "message": "个人资料更新成功",
    "user": {
      "id": 1,
      "name": "新用户名",
      "email": "新邮箱地址"
    }
  }
  ```

### 获取用户列表

- **URL**: `/users`
- **方法**: `GET`
- **请求头**: `Authorization: Bearer {token}`
- **查询参数**:
  - `page`: 页码，默认1
  - `page_size`: 每页数量，默认10
- **响应**:
  ```json
  {
    "users": [
      {
        "id": 1,
        "username": "用户名1",
        "email": "邮箱1",
        "status": "状态",
        "created_at": "创建时间"
      },
      {
        "id": 2,
        "username": "用户名2",
        "email": "邮箱2",
        "status": "状态",
        "created_at": "创建时间"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 10
  }
  ```

### 删除用户

- **URL**: `/users/{id}`
- **方法**: `DELETE`
- **请求头**: `Authorization: Bearer {token}`
- **响应**:
  ```json
  {
    "message": "用户已删除"
  }
  ```

### 更新用户状态

- **URL**: `/users/{id}/status`
- **方法**: `PUT`
- **请求头**: `Authorization: Bearer {token}`
- **请求体**:
  ```json
  {
    "status": "active|inactive|banned"
  }
  ```
- **响应**:
  ```json
  {
    "message": "用户状态已更新",
    "user": {
      "id": 1,
      "username": "用户名",
      "status": "新状态"
    }
  }
  ```

### 管理用户角色

### 管理用户角色

- **URL**: `/users/:id/roles`
- **方法**: `POST`
- **请求头**: `Authorization: Bearer {token}`
- **请求体**:
  ```json
  {
    "roles": ["user", "admin"]  // 角色名称数组
  }
  ```
- **响应**:
  ```json
  {
    "message": "角色分配成功",
    "user_id": 2,
    "roles": ["user", "admin"]
  }
  
  ```

### 获取用户角色

- **URL**: `/users/{id}/roles`
- **方法**: `GET`
- **请求头**: `Authorization: Bearer {token}`
- **响应**:
  ```json
  {
    "roles": [
      {
        "id": 1,
        "name": "admin",
        "description": "管理员"
      },
      {
        "id": 2,
        "name": "user",
        "description": "普通用户"
      }
    ]
  }
  ```

## 图片管理

### 上传图片

- **URL**: `/images/upload`
- **方法**: `POST`
- **请求头**: `Authorization: Bearer {token}`
- **表单数据**:
  - `file`: 图片文件
  - `description`: 图片描述
- **响应**:
  ```json
  {
    "message": "图片上传成功",
    "data": {
      "image_id": 1,
      "image_url": "图片URL"
    }
  }
  ```

### 批量上传图片

- **URL**: `/images/batch-upload`
- **方法**: `POST`
- **请求头**: `Authorization: Bearer {token}`
- **表单数据**:
  - `files[]`: 多个图片文件
  - `description`: 图片描述
- **响应**:
  ```json
  {
    "message": "批量上传完成",
    "results": [
      {
        "file_name": "图片1.jpg",
        "success": true,
        "image_id": 1,
        "image_url": "图片1URL"
      },
      {
        "file_name": "图片2.jpg",
        "success": false,
        "error": "错误信息"
      }
    ],
    "total": 2,
    "success_count": 1
  }
  ```

### 获取图片详情

- **URL**: `/images/{id}`
- **方法**: `GET`
- **请求头**: `Authorization: Bearer {token}`
- **响应**:
  ```json
  {
    "image": {
      "id": 1,
      "file_name": "图片.jpg",
      "description": "图片描述",
      "url": "图片URL",
      "created_at": "创建时间",
      "user_id": 1
    }
  }
  ```

### 获取图片列表

- **URL**: `/images`
- **方法**: `GET`
- **请求头**: `Authorization: Bearer {token}`
- **查询参数**:
  - `page`: 页码，默认1
  - `page_size`: 每页数量，默认10
- **响应**:
  ```json
  {
    "images": [
      {
        "id": 1,
        "file_name": "图片1.jpg",
        "description": "描述1",
        "url": "URL1",
        "created_at": "时间1"
      },
      {
        "id": 2,
        "file_name": "图片2.jpg",
        "description": "描述2",
        "url": "URL2",
        "created_at": "时间2"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 10
  }
  ```

### 搜索图片

- **URL**: `/images/search`
- **方法**: `GET`
- **请求头**: `Authorization: Bearer {token}`
- **查询参数**:
  - `keyword`: 搜索关键词
  - `page`: 页码，默认1
  - `page_size`: 每页数量，默认10
- **响应**:
  ```json
  {
    "images": [
      {
        "id": 1,
        "file_name": "图片1.jpg",
        "description": "描述1",
        "url": "URL1",
        "created_at": "时间1"
      }
    ],
    "total": 50,
    "page": 1,
    "page_size": 10
  }
  ```

### 删除图片

- **URL**: `/images/{id}`
- **方法**: `DELETE`
- **请求头**: `Authorization: Bearer {token}`
- **响应**:
  ```json
  {
    "message": "图片已删除"
  }
  ```

## 标签管理

### 获取所有标签

- **URL**: `/tags`
- **方法**: `GET`
- **请求头**: `Authorization: Bearer {token}`
- **响应**:
  ```json
  {
    "tags": [
      {
        "id": 1,
        "name": "标签1"
      },
      {
        "id": 2,
        "name": "标签2"
      }
    ]
  }
  ```

### 创建标签

- **URL**: `/tags`
- **方法**: `POST`
- **请求头**: `Authorization: Bearer {token}`
- **请求体**:
  ```json
  {
    "tag_name": "新标签"
  }
  ```
- **响应**:
  ```json
  {
    "message": "标签创建成功",
    "tag_id": 3
  }
  ```

### 为图片添加标签

- **URL**: `/tags/image/{image_id}`
- **方法**: `POST`
- **请求头**: `Authorization: Bearer {token}`
- **请求体**:
  ```json
  {
    "tag_id": 1
  }
  ```
- **响应**:
  ```json
  {
    "message": "标签添加成功"
  }
  ```

### 获取带有特定标签的图片

- **URL**: `/tags/{id}/images`
- **方法**: `GET`
- **请求头**: `Authorization: Bearer {token}`
- **查询参数**:
  - `page`: 页码，默认1
  - `page_size`: 每页数量，默认10
- **响应**:
  ```json
  {
    "images": [
      {
        "id": 1,
        "file_name": "图片1.jpg",
        "description": "描述1",
        "url": "URL1",
        "created_at": "时间1"
      }
    ],
    "total": 20,
    "page": 1,
    "page_size": 10
  }
  ```

## 私有文件管理

### 上传私有文件

- **URL**: `/private-files/upload`
- **方法**: `POST`
- **请求头**: `Authorization: Bearer {token}`
- **表单数据**:
  - `file`: 文件
  - `is_encrypted`: 是否加密 (true/false)
  - `password`: 加密密码 (如果is_encrypted为true)
- **响应**:
  ```json
  {
    "message": "文件上传成功",
    "file": {
      "id": 1,
      "file_name": "文档.pdf",
      "file_size": 12345,
      "file_type": "application/pdf",
      "is_encrypted": true,
      "created_at": "创建时间"
    }
  }
  ```

### 批量上传私有文件

- **URL**: `/private-files/batch-upload`
- **方法**: `POST`
- **请求头**: `Authorization: Bearer {token}`
- **表单数据**:
  - `files[]`: 多个文件
  - `is_encrypted`: 是否加密 (true/false)
  - `password`: 加密密码 (如果is_encrypted为true)
- **响应**:
  ```json
  {
    "message": "批量上传完成",
    "results": [
      {
        "file_name": "文档1.pdf",
        "success": true,
        "file": {
          "id": 1,
          "file_name": "文档1.pdf",
          "file_size": 12345,
          "file_type": "application/pdf",
          "is_encrypted": true,
          "created_at": "创建时间"
        }
      },
      {
        "file_name": "文档2.pdf",
        "success": false,
        "error": "错误信息"
      }
    ],
    "total": 2,
    "success_count": 1
  }
  ```

### 获取私有文件信息

- **URL**: `/private-files/{id}`
- **方法**: `GET`
- **请求头**: `Authorization: Bearer {token}`
- **查询参数**:
  - `password`: 文件密码 (如果文件已加密)
- **响应**:
  ```json
  {
    "file": {
      "id": 1,
      "file_name": "文档.pdf",
      "file_size": 12345,
      "file_type": "application/pdf",
      "is_encrypted": true,
      "created_at": "创建时间"
    }
  }
  ```

### 获取私有文件列表

- **URL**: `/private-files`
- **方法**: `GET`
- **请求头**: `Authorization: Bearer {token}`
- **查询参数**:
  - `page`: 页码，默认1
  - `page_size`: 每页数量，默认10
- **响应**:
  ```json
  {
    "files": [
      {
        "id": 1,
        "file_name": "文档1.pdf",
        "file_size": 12345,
        "file_type": "application/pdf",
        "is_encrypted": true,
        "created_at": "创建时间"
      },
      {
        "id": 2,
        "file_name": "文档2.pdf",
        "file_size": 67890,
        "file_type": "application/pdf",
        "is_encrypted": false,
        "created_at": "创建时间"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 10
  }
  ```

### 更新私有文件信息

- **URL**: `/private-files/{id}`
- **方法**: `PUT`
- **请求头**: `Authorization: Bearer {token}`
- **请求体**:
  ```json
  {
    "file_name": "新文件名.pdf",
    "is_encrypted": true,
    "password": "新密码"
  }
  ```
- **响应**:
  ```json
  {
    "message": "文件信息已更新",
    "file": {
      "id": 1,
      "file_name": "新文件名.pdf",
      "file_size": 12345,
      "file_type": "application/pdf",
      "is_encrypted": true,
      "created_at": "创建时间"
    }
  }
  ```

### 删除私有文件

- **URL**: `/private-files/{id}`
- **方法**: `DELETE`
- **请求头**: `Authorization: Bearer {token}`
- **响应**:
  ```json
  {
    "message": "文件已删除"
  }
  ```

## 令牌管理

### 创建API令牌

- **URL**: `/api/tokens`
- **方法**: `POST`
- **请求头**: `Authorization: Bearer {token}`
- **请求体**:
  ```json
  {
    "device_id": "设备ID",
    "ip_address": "IP地址"
  }
  ```
- **响应**:
  ```json
  {
    "token": {
      "token": "api_token_value",
      "expires_at": "过期时间"
    },
    "device_id": "设备ID"
  }
  ```

### 获取API令牌列表

- **URL**: `/api/tokens`
- **方法**: `GET`
- **请求头**: `Authorization: Bearer {token}`
- **响应**:
  ```json
  {
    "tokens": [
      {
        "id": 1,
        "token": "令牌1",
        "device_id": "设备ID1",
        "ip_address": "IP地址1",
        "created_at": "创建时间",
        "expires_at": "过期时间",
        "last_used_at": "最后使用时间"
      },
      {
        "id": 2,
        "token": "令牌2",
        "device_id": "设备ID2",
        "ip_address": "IP地址2",
        "created_at": "创建时间",
        "expires_at": "过期时间",
        "last_used_at": "最后使用时间"
      }
    ]
  }
  ```

### 撤销API令牌

- **URL**: `/api/tokens/{token}`
- **方法**: `DELETE`
- **请求头**: `Authorization: Bearer {token}`
- **响应**:
  ```json
  {
    "message": "token 已撤销"
  }
  ```

### 验证令牌

- **URL**: `/api/verify-token`
- **方法**: `GET`
- **请求头**: `Authorization: Bearer {token}`
- **响应**:
  ```json
  {
    "valid": true,
    "user_id": 1
  }
  ```

## 权限管理

### 获取所有权限

- **URL**: `/permissions/all`
- **方法**: `GET`
- **请求头**: `Authorization: Bearer {token}`
- **响应**:
  ```json
  {
    "permissions": [
      {
        "permission_id": 1,
        "name": "upload_img",
        "description": "上传图片权限"
      },
      {
        "permission_id": 2,
        "name": "search_img",
        "description": "搜索图片权限"
      },
      {
        "permission_id": 3,
        "name": "view_users",
        "description": "查看用户列表权限"
      }
    ]
  }
  ```

### 获取所有角色

- **URL**: `/permissions/roles`
- **方法**: `GET`
- **请求头**: `Authorization: Bearer {token}`
- **响应**:
  ```json
  {
    "roles": [
      {
        "role_id": 1,
        "role_name": "admin",
        "description": "管理员",
        "is_active": true
      },
      {
        "role_id": 2,
        "role_name": "user",
        "description": "普通用户",
        "is_active": true
      }
    ]
  }
  ```

### 获取角色的权限

- **URL**: `/permissions/roles/{role}`
- **方法**: `GET`
- **请求头**: `Authorization: Bearer {token}`
- **响应**:
  ```json
  {
    "role": "user",
    "permissions": [
      {
        "permission_id": 1,
        "name": "upload_img",
        "description": "上传图片权限"
      },
      {
        "permission_id": 2,
        "name": "search_img",
        "description": "搜索图片权限"
      },
      {
        "permission_id": 3,
        "name": "view_users",
        "description": "查看用户列表权限"
      }
    ]
  }
  ```

### 更新角色的权限

- **URL**: `/permissions/roles/{role}`
- **方法**: `PUT`
- **请求头**: `Authorization: Bearer {token}`
- **请求体**:
  ```json
  {
    "permissions": ["upload_img", "search_img", "view_users"]
  }
  ```
- **响应**:
  ```json
  {
    "message": "角色权限更新成功",
    "role": "user",
    "permissions": ["upload_img", "search_img", "view_users"]
  }
  ```

### 创建新权限

- **URL**: `/permissions/create`
- **方法**: `POST`
- **请求头**: `Authorization: Bearer {token}`
- **请求体**:
  ```json
  {
    "name": "new_permission",
    "description": "新权限描述"
  }
  ```
- **响应**:
  ```json
  {
    "message": "权限创建成功",
    "permission": "new_permission"
  }
  ```

### 创建新角色

- **URL**: `/permissions/roles/create`
- **方法**: `POST`
- **请求头**: `Authorization: Bearer {token}`
- **请求体**:
  ```json
  {
    "name": "new_role",
    "description": "新角色描述"
  }
  ```
- **响应**:
  ```json
  {
    "message": "角色创建成功",
    "role": "new_role"
  }
  ```

### 同步配置文件中的权限到数据库

- **URL**: `/permissions/sync`
- **方法**: `POST`
- **请求头**: `Authorization: Bearer {token}`
- **响应**:
  ```json
  {
    "message": "权限同步成功"
  }
  ```

### 获取指定用户的所有权限

- **URL**: `/permissions/users/{id}/permissions`
- **方法**: `GET`
- **请求头**: `Authorization: Bearer {token}`
- **响应**:
  ```json
  {
    "user_id": 2,
    "permissions": [
      {
        "permission_id": 1,
        "name": "upload_img",
        "description": "上传图片权限"
      },
      {
        "permission_id": 2,
        "name": "search_img",
        "description": "搜索图片权限"
      },
      {
        "permission_id": 5,
        "name": "view_users",
        "description": "查看用户列表权限"
      }
    ]
  }
  ```

### 获取当前用户的所有权限

- **URL**: `/permissions/users/current/permissions`
- **方法**: `GET`
- **请求头**: `Authorization: Bearer {token}`
- **响应**:
  ```json
  {
    "user_id": 1,
    "permissions": [
      {
        "permission_id": 1,
        "name": "upload_img",
        "description": "上传图片权限"
      },
      {
        "permission_id": 2,
        "name": "search_img",
        "description": "搜索图片权限"
      },
      {
        "permission_id": 3,
        "name": "createtag",
        "description": "创建标签权限"
      },
      {
        "permission_id": 4,
        "name": "manage_private_files",
        "description": "管理私有文件权限"
      },
      {
        "permission_id": 5,
        "name": "view_users",
        "description": "查看用户列表权限"
      },
      {
        "permission_id": 6,
        "name": "manage_users",
        "description": "管理用户权限"
      },
      {
        "permission_id": 7,
        "name": "manage_user_status",
        "description": "管理用户状态权限"
      },
      {
        "permission_id": 8,
        "name": "manage_user_roles",
        "description": "管理用户角色权限"
      },
      {
        "permission_id": 9,
        "name": "manage_permissions",
        "description": "管理系统权限的权限"
      }
    ]
  }
  ``` 