import { useState } from "react";

export type TimerID = number;

export type Debounce = { reset: (onComplete: () => void) => void };

export function useDebouncing(timeoutMS: number): Debounce {
  const [timer, setTimer] = useState<undefined | number>(undefined);

  return {
    reset: (onComplete: () => void) => {
      if (timer) {
        clearTimeout(timer);
      }

      const t = setTimeout(
        onComplete,
        timeoutMS,
      );
      setTimer(t);
    },
  };
}
