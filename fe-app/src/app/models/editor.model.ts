export const initialMockApiJsonString : string = [
    '{',
    '    "name": "",',
    '    "url": "", ',
    '    "responses": {',
    '        "get": {},',
    '        "delete": {},',
    '        "post": {},',
    '        "patch": {}',
    '     }',
    '}'
  ].join('\n')

export const enum  notificationLevel {
    info = 1, 
    warning = 2,
    error = 3
  }
