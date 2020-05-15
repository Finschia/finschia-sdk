package keeper

// Hooks wrapper struct for safety box keeper
type Hooks struct {
	keeper Keeper
}

// Return the wrapper struct
func (keeper Keeper) Hooks() *Hooks {
	return &Hooks{keeper}
}
