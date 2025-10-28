import type { ReactNode } from "react";

export function centerContent(content: ReactNode) {
    return (
        <div className="flex-1 flex items-center justify-center">
            {content}
        </div>
    )
}