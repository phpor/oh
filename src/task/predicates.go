// Released under an MIT-style license. See LICENSE.

// Generated by create-predicates.oh

package task

import . "github.com/michaelmacinnis/oh/src/cell"

func bind_predicates(s *Scope) {

	s.DefineMethod("is-atom", func(t *Task, args Cell) bool {
		return t.Return(NewBoolean(IsAtom(Car(args))))
	})

	s.DefineMethod("is-boolean", func(t *Task, args Cell) bool {
		return t.Return(NewBoolean(IsBoolean(Car(args))))
	})

	s.DefineMethod("is-builtin", func(t *Task, args Cell) bool {
		return t.Return(NewBoolean(IsBuiltin(Car(args))))
	})

	s.DefineMethod("is-channel", func(t *Task, args Cell) bool {
		return t.Return(NewBoolean(IsChannel(Car(args))))
	})

	s.DefineMethod("is-cons", func(t *Task, args Cell) bool {
		return t.Return(NewBoolean(IsCons(Car(args))))
	})

	s.DefineMethod("is-float", func(t *Task, args Cell) bool {
		return t.Return(NewBoolean(IsFloat(Car(args))))
	})

	s.DefineMethod("is-integer", func(t *Task, args Cell) bool {
		return t.Return(NewBoolean(IsInteger(Car(args))))
	})

	s.DefineMethod("is-method", func(t *Task, args Cell) bool {
		return t.Return(NewBoolean(IsMethod(Car(args))))
	})

	s.DefineMethod("is-null", func(t *Task, args Cell) bool {
		return t.Return(NewBoolean(IsNull(Car(args))))
	})

	s.DefineMethod("is-number", func(t *Task, args Cell) bool {
		return t.Return(NewBoolean(IsNumber(Car(args))))
	})

	s.DefineMethod("is-object", func(t *Task, args Cell) bool {
		return t.Return(NewBoolean(IsContext(Car(args))))
	})

	s.DefineMethod("is-pipe", func(t *Task, args Cell) bool {
		return t.Return(NewBoolean(IsPipe(Car(args))))
	})

	s.DefineMethod("is-rational", func(t *Task, args Cell) bool {
		return t.Return(NewBoolean(IsRational(Car(args))))
	})

	s.DefineMethod("is-status", func(t *Task, args Cell) bool {
		return t.Return(NewBoolean(IsStatus(Car(args))))
	})

	s.DefineMethod("is-string", func(t *Task, args Cell) bool {
		return t.Return(NewBoolean(IsString(Car(args))))
	})

	s.DefineMethod("is-symbol", func(t *Task, args Cell) bool {
		return t.Return(NewBoolean(IsSymbol(Car(args))))
	})

	s.DefineMethod("is-syntax", func(t *Task, args Cell) bool {
		return t.Return(NewBoolean(IsSyntax(Car(args))))
	})
}
