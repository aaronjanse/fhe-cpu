number_of_cells = int(input(''))


cells = 'cell(in_left=in_left,    in_right=cell_in_R001, out_left=out_left,   out_right=cell_in_L001, should_move=should_move, direction=direction)\n'

cell_line_format_str = '    cell(in_left=cell_in_L{0:03d}, in_right=cell_in_R{1:03d}, out_left=cell_in_R{0:03d}, out_right=cell_in_L{1:03d}, should_move=should_move, direction=direction)\n'
for i in range(1, number_of_cells-1):
    cells += cell_line_format_str.format(i, i+1)

cells += '    cell(in_left=cell_in_L{0:03d}, in_right=in_right,   out_left=cell_in_R{0:03d}, out_right=out_right,  should_move=should_move, direction=direction)'.format(number_of_cells-1)

print('''STATEFUL CHIP cell{} {{
    IN in_left, in_right, should_move, direction;
    OUT out_left, out_right;

    PARTS:
    {}
}}'''.format(number_of_cells, cells))