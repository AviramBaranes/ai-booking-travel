import { booking } from "../client";
import { withErrorHandler } from "./_api";

export function searchAvailableCars(p: booking.SearchAvailabilityRequest) {
  return withErrorHandler((client) => client.booking.SearchAvailability(p));
}

export function bookCar(p: booking.BookRequest) {
  return withErrorHandler((client) => client.booking.Book(p));
}
