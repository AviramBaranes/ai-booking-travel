import { ButtonHTMLAttributes, PropsWithChildren } from "react";

import { Loading } from "./Loading";

type ButtonProps = PropsWithChildren<
  { loading?: boolean } & ButtonHTMLAttributes<HTMLButtonElement>
>;

export function Button({ loading, children, disabled, ...props }: ButtonProps) {
  return (
    <button disabled={loading || disabled} {...props}>
      {loading ? <Loading className="mx-auto" /> : children}
    </button>
  );
}
