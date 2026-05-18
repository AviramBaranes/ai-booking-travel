import { availability, booking_handlers } from "../client";
import { withErrorHandler } from "./_api";

export function searchAvailableCars(p: availability.SearchAvailabilityParams) {
  return withErrorHandler((client) => client.booking.SearchAvailability(p));
}

export function bookCar(p: booking_handlers.BookParams) {
  return withErrorHandler((client) => client.booking.Book(p));
}
