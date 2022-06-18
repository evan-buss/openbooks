import debounce from "lodash/debounce";
import { useEffect, useState } from "react";

export function useHeight(): number {
  const [height, setHeight] = useState(window.innerHeight);

  useEffect(() => {
    const handleResize = debounce(() => setHeight(window.innerHeight), 100);

    window.addEventListener("resize", handleResize);

    return () => {
      window.removeEventListener("resize", handleResize);
    };
  }, []);

  return height;
}
