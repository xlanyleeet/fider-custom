import React from "react"
import { classSet, FiderContext } from "@fider/services"
import { Trans } from "@lingui/react/macro"

import "./PoweredByFider.scss"

interface PoweredByFiderProps {
  slot: string
  className?: string
}

export const PoweredByFider = (props: PoweredByFiderProps) => {
  const fider = React.useContext(FiderContext)
  const source = encodeURIComponent(window?.location?.host || "")
  const medium = "powered-by"
  const campaign = props.slot
  const version = fider.settings?.version?.split("-")[0]
  const versionString = fider.isSingleHostMode() && version && version !== "dev" ? `${version}` : ""
  const sourceRevision = /^gradient-([a-f0-9]{7,40})$/.exec(fider.settings?.version || "")?.[1] || "gradient"

  const className = classSet({
    "c-powered": true,
    [props.className || ""]: props.className,
  })

  return (
    <div className={className}>
      <a rel="noopener" className="text-2xs" href={`https://fider.io?utm_source=${source}&utm_medium=${medium}&utm_campaign=${campaign}`} target="_blank">
        Powered by Fider ⚡
      </a>
      {versionString && <span className="text-2xs block">{versionString}</span>}
      <a className="text-2xs block" href={`https://github.com/xlanyleeet/fider-custom/tree/${sourceRevision}`} target="_blank" rel="noopener">
        <Trans id="footer.source">Source code (AGPL-3.0)</Trans>
      </a>
    </div>
  )
}
