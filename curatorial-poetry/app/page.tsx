import Link from 'next/link'

async function getData() {
  const res = await fetch('https://v9tzav4gnk.execute-api.us-east-1.amazonaws.com/Prod/')
  // The return value is *not* serialized
  // You can return Date, Map, Set, etc.
 
  if (!res.ok) {
    // This will activate the closest `error.js` Error Boundary
    throw new Error('Failed to fetch data')
  }
 
  return res.json()
}

interface IPoem {
  id: string
  poem: string
  url: string
}

function getTaggedText(text: string, tag: string) {
  var startTag = "<" + tag + ">"
  var startTagLength = startTag.length
  var startPos = text.indexOf("<" + tag + ">")
  startPos = startPos + startTagLength
  var endPos = text.indexOf("</" + tag + ">")

  return text.substring(startPos, endPos)
}

export default async function Home() {
  const poems: IPoem[] = await getData()
  const arrayDataItems = poems.map(poem => 
    <li key={poem.id}>
      <h1><Link href={poem.url}>{getTaggedText(poem.poem, "title")}</Link></h1>
      <br/>
      <p>{getTaggedText(poem.poem, "poem")}</p>
      <br />
    </li>

  )

  return (
    <main className="flex min-h-screen flex-col items-center justify-between p-24">
      <div className="z-10 max-w-5xl w-full items-center justify-between font-mono text-sm lg:flex">

      </div>

      <div className="relative flex place-items-center before:absolute before:h-[300px] before:w-[480px]">
       <ul>{arrayDataItems}</ul>
      </div>

      <div className="mb-32 grid text-center lg:max-w-5xl lg:w-full lg:mb-0 lg:grid-cols-4 lg:text-left">

      </div>
    </main>
  )
}
