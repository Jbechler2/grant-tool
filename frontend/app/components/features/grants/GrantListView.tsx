import { Grant } from "@/types";
import { useRouter } from "next/navigation";
import GrantRow from "./GrantRow";

interface GrantListViewProps {
  ViewMode: string,
  Grants: Grant[]
}

function GrantCard({ grant }: { grant: Grant }) {
  const router = useRouter()

  return (
    <div
      onClick={() => router.push(`/grants/${grant.id}`)}
      className="bg-white border border-border rounded-lg p-5 flex flex-col gap-2 cursor-pointer hover:border-border/80 transition-colors"
    >
      <div className="flex justify-between items-center">
        <p className="text-xs text-muted-foreground">{grant.funder_name}</p>
        {grant.award_amount_max && (
          <span className="text-xs font-medium px-2 py-0.5 rounded-full bg-green-100 text-green-800">
            ${Number(grant.award_amount_max).toLocaleString()}
          </span>
        )}
      </div>
      <p className="text-sm font-medium text-foreground leading-snug">{grant.title}</p>
      {grant.description && (
        <p className="text-xs text-muted-foreground leading-relaxed line-clamp-2">
          {grant.description}
        </p>
      )}
    </div>
  )
}

export default function GrantListView(props: GrantListViewProps) {
  if(props.ViewMode === "list"){
    return (
      <div>
      <div className='grid grid-cols-4 px-4 py-2 text-sm font-medium text-muted-foreground border-b'>
        <span>Title</span>
        <span>Funder</span>
        <span>Min Amount</span>
        <span>Max Amount</span>
      </div>
        {props.Grants.map(grant => (
          <GrantRow key={grant.id} grant={grant} />
        ))}
      </div>
    )
  } else {
    return (
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
        {props.Grants.map(grant => (
          <GrantCard key={grant.id} grant={grant} />
        ))}
      </div>
    )
  }

}