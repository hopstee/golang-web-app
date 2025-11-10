import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import type { Field } from "@/types/entities";

interface SelectFieldProps {
    field: Field;
    value: string;
    onChange: (name: string, value: string) => void;
}

export default function SelectField(props: SelectFieldProps) {
    return (
        <Select value={props.value} onValueChange={(value) => props.onChange(props.field.name, value)}>
            <SelectTrigger className="w-full">
                <SelectValue placeholder={props.field.label} />
            </SelectTrigger>
            <SelectContent>
                {props.field.options?.map(option => (
                    <SelectItem value={option}>{option}</SelectItem>
                ))}
            </SelectContent>
        </Select>
    );
}