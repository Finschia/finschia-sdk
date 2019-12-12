# Concepts

## Identity and Access Management 

IAM(Identity and Access Management) manages roles and access privileges of individual users(accounts) to specific resources. 
The account may be a fungible token issuer or owner of the non-fungible token.

### States
IAM manages Triad; Account, Resource, and Action. It means that "An Account has permission to do Action for Resource".
The address of the account is the key and the list of the pair of resource and action is the value.

```
+-----------------+      +----------------+                    
|                 |  +--->Resource, Action|                    
|Account          |  |   +----------------+                    
+-----------------+  |   +----------------+                    
         |           +-->|Resource, Action|                    
+--------v--------+  |   +----------------+                    
|                 |  |   +----------------+                    
|AccountPermission|----->|Resource, Action|                    
+-----------------+      +----------------+                    
```

### Define Permission
The permission is an interface. It is defined by the module which uses this IAM module.
```go
type PermissionI interface {
	GetResource() string
	GetAction() string
	Equal(string, string) bool
}
```


### Group of Accounts
A permission of an account can be shared by others.
```
+-----------------+          +-----------------+      +----------------+                    
|                 |          |                 |  +--->Resource, Action|                    
|Account          |          |Account          |  |   +----------------+                    
+-----------------+          +-----------------+  |   +----------------+                    
         |                            |           +-->|Resource, Action|                    
+--------v--------+          +--------v--------+  |   +----------------+                    
|Inherited        |inherit   |                 |  |   +----------------+                    
|AccountPermission----------->AccountPermission|----->|Resource, Action|                    
+-----------------+          +-----------------+      +----------------+                    
```

The parent account permission can be shared by its children. It effect like a group of account permission
```
+-----------------+      +-----------------+          +-----------------+      +----------------+       
|                 |      |Inherited        |          |                 |  +--->Resource, Action|       
|Account          |------>AccountPermission|-----+    |Account          |  |   +----------------+       
+-----------------+      +-----------------+     |    +-----------------+  |   +----------------+       
                                                 |             |           +-->|Resource, Action|       
+-----------------+      +-----------------+     |    +--------v--------+  |   +----------------+       
|                 |      |Inherited        |inherit   |                 |  |   +----------------+       
|Account          |------>AccountPermission----------->AccountPermission|----->|Resource, Action|       
+-----------------+      +-----------------+     |    +-----------------+      +----------------+       
+-----------------+      +-----------------+     |                                                      
|                 |      |Inherited        |     |                                                      
|Account          |------>AccountPermission|-----+                                                      
+-----------------+      +-----------------+                                                            
```

### Independent permission to modules
A module refers to the store of this IAM module with a prefix. It enables independent account permission to each module.
Though, it is not mandatory. The module developer can decide whether to use the prefix or not.