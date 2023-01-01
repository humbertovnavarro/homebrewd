interface EditableListItemProps {
    default: string
    key: number
}
export const EditableListItem: React.FC<EditableListItemProps> = (props: EditableListItemProps) => {
    return (
        <li key={props.key}>
            <input type="text" name={props.default} defaultValue={props.default} className="bg-white/10 p-1 text-white hover:bg-white/20 rounded-md shadow-md"/>
        </li>
    )
}
