package rbac

import (
	"log"
	"strings"
	"sync"
)

// RbacCache определяет операции для обновления и инвалидации кеша RBAC
type IRbacCache interface {
	ReloadAll() error

	InvalidateRole(roleID string)
	InvalidateGrant(grantID string)
	InvalidateRoleGrant(roleGrantID string)
	InvalidateUserRole(userRoleID string)

	HasPermission(userID, code string) bool
}

// RbacCacheImpl — in-memory реализация
type RbacCacheImpl struct {
	mu sync.RWMutex

	// Maps for fast lookup
	users      map[string][]string // userID -> roleIDs
	roles      map[string]string   // roleID -> roleName
	grants     map[string]string   // grantID -> grantCode
	roleGrants map[string][]string // roleID -> grantIDs
	userRoles  map[string][]string // userRoleID -> [userID, roleID]

	logger *log.Logger
}

// NewRBACCache создает новый экземпляр кэша
func NewRBACCache(logger *log.Logger) IRbacCache {
	return &RbacCacheImpl{
		logger: logger,

		users:      make(map[string][]string),
		roles:      make(map[string]string),
		grants:     make(map[string]string),
		roleGrants: make(map[string][]string),
		userRoles:  make(map[string][]string),
	}
}

// ReloadAll — полная перезагрузка из БД (здесь подключить sqlc или GORM)
func (rcv *RbacCacheImpl) ReloadAll() error {
	rcv.mu.Lock()
	defer rcv.mu.Unlock()

	rcv.logger.Println("RbacCache: ReloadAll called — full refresh")

	// TODO: load users, roles, grants, role_grants, user_roles from DB

	return nil
}

// Инвалидация отдельных сущностей
func (rcv *RbacCacheImpl) InvalidateRole(roleID string) {
	rcv.mu.Lock()
	defer rcv.mu.Unlock()

	delete(rcv.roles, roleID)
	delete(rcv.roleGrants, roleID)

	rcv.logger.Printf("RbacCache: role %s invalidated\n", roleID)
}

func (rcv *RbacCacheImpl) InvalidateGrant(grantID string) {
	rcv.mu.Lock()
	defer rcv.mu.Unlock()
	delete(rcv.grants, grantID)
	for r, gs := range rcv.roleGrants {
		newGs := []string{}

		for _, g := range gs {
			if g != grantID {
				newGs = append(newGs, g)
			}
		}

		rcv.roleGrants[r] = newGs
	}

	rcv.logger.Printf("RbacCache: grant %s invalidated\n", grantID)
}

func (rcv *RbacCacheImpl) InvalidateRoleGrant(roleGrantID string) {
	rcv.mu.Lock()
	defer rcv.mu.Unlock()

	for r, gs := range rcv.roleGrants {
		newGs := []string{}

		for _, g := range gs {
			if g != roleGrantID {
				newGs = append(newGs, g)
			}
		}

		rcv.roleGrants[r] = newGs
	}
	rcv.logger.Printf("RbacCache: role_grant %s invalidated\n", roleGrantID)
}

func (rcv *RbacCacheImpl) InvalidateUserRole(userRoleID string) {
	rcv.mu.Lock()
	defer rcv.mu.Unlock()

	delete(rcv.userRoles, userRoleID)

	rcv.logger.Printf("RbacCache: user_role %s invalidated\n", userRoleID)
}

// HasPermission проверяет, есть ли у пользователя право на grantCode
func (rcv *RbacCacheImpl) HasPermission(userID, code string) bool {
	rcv.mu.RLock()
	defer rcv.mu.RUnlock()

	roleIDs, ok := rcv.users[userID]
	if !ok {
		return false
	}

	for _, roleID := range roleIDs {
		grantIDs := rcv.roleGrants[roleID]

		for _, grantID := range grantIDs {
			grantCode := rcv.grants[grantID]
			if matchPattern(code, grantCode) {
				return true
			}
		}
	}

	return false
}

// matchPattern — поддержка wildcard '*'
func matchPattern(code, pattern string) bool {
	if code == pattern {
		return true
	}

	codeParts := strings.Split(code, ":")
	patternParts := strings.Split(pattern, ":")

	for i := 0; i < len(patternParts) && i < len(codeParts); i++ {
		if patternParts[i] == "*" {
			continue
		}

		if patternParts[i] != codeParts[i] {
			return false
		}
	}

	return len(codeParts) == len(patternParts)
}
