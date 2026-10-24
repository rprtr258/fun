package json

func (decoder Decoder[T]) Map[R any](f func(T) R) Decoder[R] {
	return func(v any, res *R) error {
		var dest0 T
		if err := decoder(v, &dest0); err != nil {
			return nil
		}

		*res = f(dest0)
		return nil
	}
}

func Map2[T0, T1, T any](
	combine func(T0, T1) T,
	d0 Decoder[T0],
	d1 Decoder[T1],
) Decoder[T] {
	return func(v any, res *T) error {
		var dest0 T0
		if err := d0(v, &dest0); err != nil {
			return err
		}

		var dest1 T1
		if err := d1(v, &dest1); err != nil {
			return err
		}

		*res = combine(dest0, dest1)
		return nil
	}
}

func Map3[T0, T1, T2, T any](
	combine func(T0, T1, T2) T,
	d0 Decoder[T0],
	d1 Decoder[T1],
	d2 Decoder[T2],
) Decoder[T] {
	return func(v any, res *T) error {
		var dest0 T0
		if err := d0(v, &dest0); err != nil {
			return err
		}
		var dest1 T1
		if err := d1(v, &dest1); err != nil {
			return err
		}
		var dest2 T2
		if err := d2(v, &dest2); err != nil {
			return err
		}

		*res = combine(dest0, dest1, dest2)
		return nil
	}
}

func Map4[T0, T1, T2, T3, T any](
	combine func(T0, T1, T2, T3) T,
	d0 Decoder[T0],
	d1 Decoder[T1],
	d2 Decoder[T2],
	d3 Decoder[T3],
) Decoder[T] {
	return func(v any, res *T) error {
		var dest0 T0
		if err := d0(v, &dest0); err != nil {
			return err
		}
		var dest1 T1
		if err := d1(v, &dest1); err != nil {
			return err
		}
		var dest2 T2
		if err := d2(v, &dest2); err != nil {
			return err
		}
		var dest3 T3
		if err := d3(v, &dest3); err != nil {
			return err
		}

		*res = combine(dest0, dest1, dest2, dest3)
		return nil
	}
}

func Map5[T0, T1, T2, T3, T4, T any](
	combine func(T0, T1, T2, T3, T4) T,
	d0 Decoder[T0],
	d1 Decoder[T1],
	d2 Decoder[T2],
	d3 Decoder[T3],
	d4 Decoder[T4],
) Decoder[T] {
	return func(v any, res *T) error {
		var dest0 T0
		if err := d0(v, &dest0); err != nil {
			return err
		}
		var dest1 T1
		if err := d1(v, &dest1); err != nil {
			return err
		}
		var dest2 T2
		if err := d2(v, &dest2); err != nil {
			return err
		}
		var dest3 T3
		if err := d3(v, &dest3); err != nil {
			return err
		}
		var dest4 T4
		if err := d4(v, &dest4); err != nil {
			return err
		}

		*res = combine(dest0, dest1, dest2, dest3, dest4)
		return nil
	}
}

func Map8[T0, T1, T2, T3, T4, T5, T6, T7, T any](
	combine func(T0, T1, T2, T3, T4, T5, T6, T7) T,
	d0 Decoder[T0],
	d1 Decoder[T1],
	d2 Decoder[T2],
	d3 Decoder[T3],
	d4 Decoder[T4],
	d5 Decoder[T5],
	d6 Decoder[T6],
	d7 Decoder[T7],
) Decoder[T] {
	return func(v any, res *T) error {
		var dest0 T0
		if err := d0(v, &dest0); err != nil {
			return err
		}
		var dest1 T1
		if err := d1(v, &dest1); err != nil {
			return err
		}
		var dest2 T2
		if err := d2(v, &dest2); err != nil {
			return err
		}
		var dest3 T3
		if err := d3(v, &dest3); err != nil {
			return err
		}
		var dest4 T4
		if err := d4(v, &dest4); err != nil {
			return err
		}
		var dest5 T5
		if err := d5(v, &dest5); err != nil {
			return err
		}
		var dest6 T6
		if err := d6(v, &dest6); err != nil {
			return err
		}
		var dest7 T7
		if err := d7(v, &dest7); err != nil {
			return err
		}

		*res = combine(dest0, dest1, dest2, dest3, dest4, dest5, dest6, dest7)
		return nil
	}
}
