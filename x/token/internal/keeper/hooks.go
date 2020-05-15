package keeper

// Hooks wrapper struct for safety box keeper
type Hooks struct {
	k Keeper
}

// Return the wrapper struct
func (k Keeper) Hooks() *Hooks {
	return &Hooks{k}
}
