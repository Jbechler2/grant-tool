import Link from "next/link"
import { Card, CardTitle, CardContent } from "@/components/ui/card"
import { Grant } from "@/types"

export default function ApplicationCard(grant: Grant) {
  return (
  <Link 
    key={grant.id}
    href=""
  >
    <Card            
      className="max-w-sm hover:shadow-md"
    >
      <CardTitle>{grant.title}</CardTitle>
    </Card>
  </Link>
)

}
