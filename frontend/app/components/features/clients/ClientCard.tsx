import Link from "next/link"
import { Card, CardTitle } from "@/components/ui/card"
import { Client } from "@/types"

export default function ClientCard(client: Client) {
  return (
  <Link 
    key={client.id}
    href=""
  >
    <Card            
      className="max-w-sm hover:shadow-md"
    >
      <CardTitle>{client.name}</CardTitle>
    </Card>
  </Link>
)

}
