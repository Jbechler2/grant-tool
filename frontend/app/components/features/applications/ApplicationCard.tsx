import Link from "next/link"
import { Card, CardTitle, CardContent } from "@/components/ui/card"
import { Application } from "@/types"
import { Calendar } from 'lucide-react'

export default function ApplicationCard(app: Application) {
  return (
  <Link 
    key={app.id}
    href=""
    className="flex"
  >
    <div className="bg-white rounded-lg shadow-xl pr-6 pb-6 flex flex-col justify-between lg:h-40 md:h-32 sm:h-20 font-halant">
      <div className="flex flex-row justify-between pt-4 pl-4">
        <div className="grid place-items-center bg-seafoam mx-auto my-auto p-2 rounded-md border-2 border-forest h-10 w-10">
          {getIconLetters(app.title)}
        </div>
        <div className="px-3">
          <div>
            <h1 className="lg:text-lg font-bold  pb-0 mb-0">{app.title}</h1>
          </div>
          <div className="grid place-items-start p-0 m-0">
            <h3 className="pt-0 mt-0 italic">{app.status}</h3>
          </div>
        </div>
        <div>
          <h1 className="lg:text-lg font-extrabold">$100,000</h1>
        </div>
      </div>
      <hr></hr>
      <div className="flex flex-col justify-between flex-1 pl-4 pt-2">
        <div>
          <p className="whitespace-normal break-words"></p>
        </div>
        <div>Topics</div>
      </div>
      <hr></hr>
      <div className="flex flex-row">
        <div className="flex flex-row">
          <div className="grid place-items-center">
            <Calendar className="pl-4 w-15"/> 
          </div>
          <div>
            <div>Deadline</div>
            <div>Mar 31, 2026</div>
          </div>
        </div>
        <div>
          Details
        </div>
      </div>
    </div>
  </Link>
)

}

function getIconLetters(name: string): string{
  const words = name.split(' ')
  if(words.length > 0){
    if(words.length == 1){
    return words[0][0].toUpperCase()
  } if(words.length == 2){
    return words[0][0].toUpperCase() + words[1][0].toUpperCase();
  } else {
    const letter1 = words[0][0].toUpperCase()
    let letter2 = ''
    words.sort((a, b) => {
      return b.length - a.length
    })
    if(words[0][0].toUpperCase() === letter1){
      letter2 = words[1][0].toUpperCase()
    } else {
      letter2 = words[0][0].toUpperCase()
    }
    return letter1 + letter2
  }
  } else {
    return 'AA'
  }
  
}