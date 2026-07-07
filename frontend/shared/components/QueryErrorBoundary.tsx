"use client";

import { Component, ReactNode } from "react";

type QueryErrorBoundaryProps = {
  children: ReactNode;
  fallback: ReactNode;
};

type QueryErrorBoundaryState = {
  hasError: boolean;
};

/**
 * Catches errors thrown during render (including thrown `useSuspenseQuery`
 * errors) and renders a fallback instead of letting the error bubble up to the
 * route boundary and re-trigger the suspense query in a loop.
 */
export class QueryErrorBoundary extends Component<
  QueryErrorBoundaryProps,
  QueryErrorBoundaryState
> {
  state: QueryErrorBoundaryState = { hasError: false };

  static getDerivedStateFromError(): QueryErrorBoundaryState {
    return { hasError: true };
  }

  render() {
    if (this.state.hasError) {
      return this.props.fallback;
    }
    return this.props.children;
  }
}
