import { Client } from "@/types";


export default function ClientRow({ client }: {client : Client}) {
  return (
    <div className="flex flex-row bg-gray-300 rounded-lg border-1 p-5">
      <div>{client.name}</div>
      <div>{client.contact_phone}</div>
    </div>
  )
}