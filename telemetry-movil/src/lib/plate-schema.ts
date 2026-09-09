import { z } from 'zod';

// plateSchema validates the vehicle plate format ABC-123 (three uppercase letters,
// a dash, three digits). Matches the regex enforced by vehicle-service (SPEC 02).
export const plateSchema = z
  .string()
  .regex(/^[A-Z]{3}-[0-9]{3}$/, 'Invalid plate format, expected ABC-123');

// plateFormSchema is the object form for React Hook Form (form values: { plate: string }).
export const plateFormSchema = z.object({ plate: plateSchema });
