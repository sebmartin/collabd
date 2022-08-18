package models

// An interface for the code logic for each game's stage. These can be thought of
// as a "state" in the game's rules-based finite state machine. The game engine will
// execute the `Run()` method once the game enters a stage.
//
// There should only ever be a single `StageRunner` running at once per game session.
// However, the `Run()` method is invoked from a go subroutine so thread safety should
// be considered when accessing shared resources.
//
// The runner receives player events via the channel passed into the `Run()` method
// argument. It should then act on each event accoringly. The runner can delegate
// the game execution to another stage by ending the `Run()` and returning the
// next `StageRunner` that will take over.
type StageRunner interface {
	Run(<-chan PlayerEvent) StageRunner
}
