import { useState } from "react"

// export function Page({ analyzedData }) {
//     return (
//         <div className="page">
//             {analyzedData.map((para, paraIdx) => {
//                 <p key={paraIdx}>
//                     {para.sentences.map((sentence, sentIdx) => (
//                         <Sentence 
//                             key={sentIdx} 
//                             sentence={sentence}
//                         />
//                     ))}
//                 </p>
//             })}
//         </div>
//     )
// }

export function Page({ analyzedData }) {
    return (
      <div className="page">
        {analyzedData.map((para, paraIdx) => (
          <p key={paraIdx}>
            {para.sentences.map((sentence, sentIdx) => (
              <span key={sentIdx}>
                <Sentence key={sentIdx} sentence={sentence} />
              </span>
            ))}
          </p>
        ))}
      </div>
    );
  }

function Sentence({ sentence }) {
    const [showFix, setShowFix] = useState(false);



    const hasIssue = sentence.isPassiveVoice && sentence.potentialFix;

    return (
        <>
            <span className={hasIssue && !showFix ? 'sentence has-issue' : showFix ? 'sentence sentence-fixed' : 'sentence'}
            onClick={() => {
                hasIssue && setShowFix(!showFix)
                
            }}
            >
                {showFix && sentence.potentialFix ? sentence.potentialFix : sentence.content}
            </span>
            {' '}
        </>
    )
}