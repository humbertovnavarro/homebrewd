interface ListItemProps {
    children: React.ReactNode
    key: number
}
export const ListItem: React.FC<ListItemProps> = (props: ListItemProps) => {
    return (
        <li key={props.key} className="bg-white/10 flex p-1 text-white hover:bg-white/20 rounded-md shadow-md w-full">
            {props.children}
        </li>
    )
}
