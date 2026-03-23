import { Grant } from "@/types";
import { cookies } from "next/headers";

export async function getGrant(id: string): Promise<Grant> {
  const token = (await cookies()).get('token')?.value;

  const url = `${process.env.NEXT_PUBLIC_API_URL}/api/v1/grants/${id}`
  const res = await fetch(url, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  if (!res.ok) throw new Error('Failed to fetch grant');
  return res.json();
}